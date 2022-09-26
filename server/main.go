package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"stl/model"
	"stl/pages"
	"stl/server/useragent"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	pgsqlDB *sql.DB

	dbInfo  = flag.String("db-info", "postgres://postgres:@localhost/savethislife?sslmode=disable", "database connection string")
	sslcert = flag.String("sslcert", "", "path to an SSL cert file")
	sslkey  = flag.String("sslkey", "", "path to an SSL key file")
	port    = flag.Int("port", 8080, "port to be sued by the server")
	host    = flag.String("host", "localhost:8080", "host used to reach this server")

	views     = flag.String("views", "views", "path to views folder")
	public    = flag.String("public", "public", "path to public folder")
	imagePets = flag.String("image-pets", "/tmp/", "path to folder containing image of pets")
	sitemaps  = flag.String("sitemaps", "", "path to folder containing sitemaps")
)

func main() {
	flag.Parse()

	if _, err := os.Stat(filepath.Join(*views)); os.IsNotExist(err) {
		log.Fatalf("views directory %q not exist!", *views)
		return
	}

	model.DBInfo = *dbInfo
	err := model.InitDB()
	if err != nil {
		log.Fatalf("cannot initialize db: %v", err)
		return
	}

	err = parseTemplates(*views)
	if err != nil {
		log.Fatalf("error when parsing views templates: %v", err)
		return
	}

	pages.PublicFolder = *public
	// if *sslcert != "" && *sslkey != "" {
	// 	model.AppURL = fmt.Sprintf("https://%v", *host)
	// } else {
	// 	model.AppURL = fmt.Sprintf("http://%v", *host)
	// }
	model.AppURL = fmt.Sprintf("https://%v", *host)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(*public))))
	router.PathPrefix("/image-pets/").Handler(http.StripPrefix("/image-pets/", http.FileServer(http.Dir(*imagePets))))
	router.PathPrefix("/sitemap/").Handler(http.StripPrefix("/sitemap/", http.FileServer(http.Dir(*sitemaps))))
	router.PathPrefix("/new-sitemaps/").Handler(http.StripPrefix("/new-sitemaps/", http.FileServer(http.Dir(*sitemaps))))

	router.HandleFunc("/", pages.Index).Methods("GET")
	// WORKAROUND the problem in .htaccess
	router.HandleFunc("/index.html.var", pages.Index).Methods("GET")

	router.HandleFunc("/robots.txt", pages.RobotsText).Methods("GET")
	router.HandleFunc("/about-us", pages.AboutUs).Methods("GET")
	router.HandleFunc("/about-us/package-info", pages.PackageInfo).Methods("GET")
	router.HandleFunc("/microchip", pages.Microchip).Methods("GET")
	router.HandleFunc("/pet-health-insurance", pages.PetHealthInsurance).Methods("GET")
	router.HandleFunc("/frequently-asked-questions", pages.Faq).Methods("GET")
	router.HandleFunc("/order-received", pages.OrderReceived).Methods("GET")
	router.HandleFunc("/geoip", pages.Geoip).Methods("GET")
	router.HandleFunc("/recents", pages.Recents).Methods("GET")
	router.HandleFunc("/lost-found-pet", pages.LostFoundPet).Methods("GET")
	router.HandleFunc("/privacy-policy", pages.PrivacyPolicy).Methods("GET")
	router.HandleFunc("/registration-types", pages.RegistrationTypes).Methods("GET")

	router.HandleFunc("/sitemap_index.xml", pages.SitemapIndex).Methods("GET")
	router.HandleFunc("/sitemap-part/{part}", pages.SitemapPart).Methods("GET")
	router.HandleFunc("/sitemap.xsl", pages.SitemapXSL).Methods("GET")

	router.HandleFunc("/contact-us", pages.GetContactUs).Methods("GET")
	router.HandleFunc("/contact-us", pages.PostContactUs).Methods("POST")

	router.HandleFunc("/contact-save-this-life", pages.GetContactStl).Methods("GET")
	router.HandleFunc("/contact-save-this-life", pages.PostContactStl).Methods("POST")

	router.HandleFunc("/shopping-microchip", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://po.savethislife.com/signup", http.StatusPermanentRedirect)
	}).Methods("GET")
	// router.HandleFunc("/shopping-microchip", pages.PostShoppingMicrochip).Methods("POST")

	router.HandleFunc("/register-stl-microchip", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://po.savethislife.com/signup", http.StatusPermanentRedirect)
	}).Methods("GET")
	// router.HandleFunc("/register-stl-microchip", pages.PostRegisterStlMicrochip).Methods("POST")

	router.HandleFunc("/searchmicrochip", pages.GetSearchMicrochip).Methods("GET")
	router.HandleFunc("/why-did-you-savethislife", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.savethislife.com", http.StatusPermanentRedirect)
	}).Methods("GET")

	router.HandleFunc("/petregistration", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://po.savethislife.com/signup", http.StatusPermanentRedirect)
	}).Methods("GET")

	router.HandleFunc("/help", pages.Help).Methods("GET")
	router.HandleFunc("/found", pages.HandleFoundPetGET).Methods("GET")

	router.HandleFunc("/f1", func(w http.ResponseWriter, r *http.Request) {
		// Request To Transfer Ownership of a Pet to Me
		http.Redirect(w, r, "https://zfrmz.com/lto2Yt676f3EOHBuJ9jL", http.StatusTemporaryRedirect)
	}).Methods("GET")

	for _, r := range httpGoneRoutes() {
		router.HandleFunc(r, pages.HttpGone).Methods("GET")
	}

	goneURLs, err := model.GoneURLs()
	if err != nil {
		log.Fatalf("error when getting gone urls: %v", err)
		return
	}
	for _, r := range goneURLs {
		router.HandleFunc(r, pages.HttpGone).Methods("GET")
	}

	router.Handle("/{microchip}",
		useragent.UserAgentMiddleware(http.HandlerFunc(pages.GetMicrochip))).
		Methods("GET")
	router.HandleFunc("/{microchip}", pages.PostMicrochip).Methods("POST")

	if *sslcert != "" && *sslkey != "" {
		// redirect every http request to https
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
		log.Printf("server listening at :%v on https...", *port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", *port),
			*sslcert, *sslkey, handlers.LoggingHandler(os.Stdout, CaselessMatcher(router))))
	} else {
		log.Printf("server listening at :%v on http...", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port),
			handlers.LoggingHandler(os.Stdout, CaselessMatcher(router))))
	}
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// remove/add not default ports from req.Host
		if req.Host == "104.198.22.132" || req.Host == "savethislife.com" {
			req.URL.Host = "www.savethislife.com"
			http.Redirect(w, req, req.URL.String(), http.StatusPermanentRedirect)
			return
		}

		req.URL.Path = strings.ToLower(req.URL.Path)
		next.ServeHTTP(w, req)
	})
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusPermanentRedirect)
}

func parseTemplates(dir string) error {
	var allFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := file.Name()
		allFiles = append(allFiles, filepath.Join(dir, filename))
	}

	t := template.New("").Funcs(template.FuncMap{
		"set": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
			renderArgs[key] = value
			return template.JS("")
		},
		"year": func() int {
			return time.Now().Year()
		},
	})

	// parses all .html files in the 'dir' folder
	templates, err := t.ParseFiles(allFiles...)
	if err != nil {
		return err
	}

	pages.Templates = templates

	return nil
}

func httpGoneRoutes() []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(goneRoutes))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

const (
	goneRoutes = `
/90016400073262
/newstl/amitriptyline/lorazepam-vs-alprazolam-treatment-of-panic-disorder.php
/99354794
/900164000043738-3
/20001105369
/991001000414409c
/90016400060531
/73024317
/98100000328
/20001061816
/90016400075724
/84519593
/newstl/gabapentin/valium-et-seresta.php
/04617541
/04306026
/newstl/gabapentin/can-i-take-adipex-if-i-have-high-blood-pressure.php
/20003512547
/99100100033581
/91001000602462
/41326333
/23094612
/991001000233097v
/20001715924
/4.8097E+250
/9910010004644
/039322555-2
/products-page/checkout/
/9851210054266v
/demo2/mike-arms-microchip-movement-2013/
/03366029
/35119379
/99100100446750
/90016400106807
/99100100061051
/900164000703268x%60
/90004000502244
/005535885/feed
/900164000043714-3
/newstl/amitriptyline/xanax-4mg-bars.php
/57374119
/900164000290326c
/9910010002467076
/62038302
/900164000405086f
/20001123691
/900164000994291v
/03372865
/wp-content/sedlex/translations/nabumetone/diabetes-in-arabic-health.php
/15126857
/900164000987551c
/991001000186779c
/99100100575268
/99100100063779
/wp-content/sedlex/backup-scheduler/map.php
/900164000597456v
/51575607
/03527581
/4043073
/newstl/prevacid/tramadol-vs-ponstan.php
/newstl/sinemet/soma-akam-ezan-saati.php
/9851210091212v
/57528094
/20005124214
/06739990
/90016400029326
/991001000263191/feed
/57348807
/13329848
/is011784
/900164000525962v
/97523106
/3353793
/48530639
/900164000739545c
/newstl/amitriptyline/smoking-weed-mixed-with-xanax.php
/0900164000209185
/07063125
/460B0C5C
/20003490523
/5601111
/08875818
/microchip/register-any-brand-of-microchip/
/tag/71270856/
/900164001128198%C2%A0
/enigma-pervice/register-any-brand-of-microchip
/900164001091020%C2%A0
/900164000760913%C2%A0
/900164000923594%C2%A0
/900164001091015%C2%A0
/register-any-brand-of-microchip
/author/stl_admin/page/31/
/author/stl_admin/page/7209/
/author/stl_admin/page/26/
/author/stl_admin/page/10968/
/author/stl_admin/page/3748/
/author/stl_admin/page/5473/
/author/stl_admin/page/22/
/category/tags/page/35/
/category/tags/page/5927/
/author/stl_admin/page/11798/
/category/tags/page/2311/
/author/stl_admin/page/5693/
/AVID009-294-777
/category/tags/page/3337/
/category/tags/page/17/
/category/tags/page/26/
/author/stl_admin/page/7054/
/category/tags/page/5/
/author/stl_admin/page/2913/
/fdx-a0a120a414
/author/stl_admin/page/6939/
/author/stl_admin/page/10694/
/author/stl_admin/page/7230/
/author/stl_admin/page/3094/
/author/stl_admin/page/7055/
/sitemap/part2849.xml
/sitemap/part2112.xml
/page/2911
/AVID078633825
/author/stl_admin/page/10870/
/author/stl_admin/page/7681/
/AVID076318266
/category/tags/page/9375/
/sitemap/part111.xml
/category/tags/page/21145/
/author/stl_admin/page/9/
/author/stl_admin/page/7675/
/category/tags/page/13336/
/wp-content/sedlex/inline_scripts/map.php
/category/tags/page/3614/
/sitemap/part1341.xml
/wp-content/sedlex/inline_styles/map.php
/category/tags/page/14/
/sitemap/part2376.xml
/0000/page/12881
/sitemap/part2939.xml
/sitemap/part2716.xml
/author/stl_admin/page/7/
/author/stl_admin/page/25/
/wp-content/sedlex/translations/map.php
/rgister-from-biz
/900164000734695/page/1782
/author/stl_admin/page/8162/
/category/tags/page/6104/
/author/stl_admin/page/1898/
/0000/page/15002
/author/stl_admin/page/5126/
/author/stl_admin/page/3473/
/category/tags/page/33/
/category/tags/page/9160/
/900164001071224/feed
/index.php
/900164001065796/feed
/900164000701585/feed
/pricing-and-prizing
/991001000244838/feed
/991001000333858/feed
/900164000571319/feed
/991001000263191/feed
/991001000233097v
/57348807
/460B0C5C
/900164000739545c
/90004000502244
/23094612
/newstl/amitriptyline/xanax-4mg-bars.php
/99100100061051
/pricing-and-prizing
/900164000571319/feed
/newstl/sinemet/soma-akam-ezan-saati.php
/90016400106807
/48530639
/newstl/gabapentin/valium-et-seresta.php
/20001105369
/991001000241197991001000241197
/900164000290326c
/991001000333858/feed
/900164000043714-3
/932001000579647932001000579647
/newstl/gabapentin/can-i-take-adipex-if-i-have-high-blood-pressure.php
/57528094
/0900164000209185
/900164000703268x%60
/newstl/amitriptyline/lorazepam-vs-alprazolam-treatment-of-panic-disorder.php
/900164000701585/feed
/99100100575268
/20001061816
/900164000994291v
/demo2/mike-arms-microchip-movement-2013/
/900164000901289900164000901289
/04617541
/5601111
/newstl/prevacid/tramadol-vs-ponstan.php
/900164000352284900164000352284
/99354794
/4043073
/20003512547
/99100100063779
/900164000399538900164000399538
/90016400029326
/005535885/feed
/62038302
/57374119
/20003490523
/9910010004644
/991001001123006991001001123006
/13329848
/products-page/checkout/
/20001715924
/900164000415745900164000415743
/07063125
/20001123691
/982000153608122982000153608122
/991001000333908991001000333908
/97523106
/900164000042053-3
/90016400062990
/90016400075724
/03372865
/03366029
/newstl/amitriptyline/soma-350-mg-get-high.php
/900164000043738-3
/wp-content/sedlex/translations/nabumetone/diabetes-in-arabic-health.php
/9851210091212v
/9851210054266v
/90016400060531
/wp-content/sedlex/backup-scheduler/map.php
/900164000525962v
/03527581
/04306026
/08875818
/98100000328
/039322555-2
/900164000991756900164000991756
/20005124214
/84519593
/15126857
/991001000186779c
/51575607
/99100100033581
/73024317
/35119379
/9910010002467076
/900164001065796/feed
/900164001071224/feed
/newstl/amitriptyline/smoking-weed-mixed-with-xanax.php
/900164000987551c
/900164000405086f
/91001000602462
/3353793
/41326333
/90016400073262
/900164000597456v
/991001000244838/feed
/991001000414409c
/06739990
/99100100446750
/0000/page/26854
/newstl/amitriptyline/smoking-weed-mixed-with-xanax.php
/newstl/gabapentin/can-i-take-adipex-if-i-have-high-blood-pressure.php
/57374119
/57528094
/04617541
/900164000040458-2
/page/11954
/category/tags/page/25/
/newstl/sinemet/soma-akam-ezan-saati.php
/900164000739545c
/08875818
/page/2060
/page/666
/wp-content/sedlex/translations/nabumetone/diabetes-in-arabic-health.php
/62038302
/3353793
/97523106
/9910010002467076
/products-page/your-account/
/page/14609
/products-page/checkout/
/991001000233097v
/84519593
/5601111
/900164000040685-2
/page/23101
/900164000571319/feed
/900164000043738-3
/99100100575268
/page/8799
/900164000405086f
/900164000994291v
/51575607
/985121005306048-2
/0000/page/10712
/900164000701585/feed
/9851210054266v
/57348807
/99354794
/author/stl_admin/page/80/
/0000/page/15492
/pricing-and-prizing
/90016400029326
/20001105369
/90016400060531
/991001000186779c
/page/31143
/71796010
/48530639
/900164000290326c
/page/1429
/is011784
/4043073
/900164000038911-2
/900164000041919-2
/900164000043770-2
/0000/page/15892
/0000/page/8304
/900164001071224/feed
/039322555-2
/991001000263191/feed
/90016400073262
/35119379
/98100000328
/900164000043423-2
/900164001128198%C2%A0
/9851210091212v
/4.8097E+250
/99100100033581
/13329848
/0000/page/11855
/20005124214
/9910010004644
/99100100446750
/0a111d5a76-2
/page/19809
/20003490523
/newstl/gabapentin/valium-et-seresta.php
/900164000525962v
/23094612
/tag/900164001136429
/reregister-2/goggles/
/0000/page/24738
/register-any-brand-of-microchip
/newstl/prevacid/tramadol-vs-ponstan.php
/900164000040042-3
/page/18982
/460B0C5C
/20001123691
/99100100061051
/91001000602462
/rgister-from-biz
/90016400075724
/900164000597456v
/03372865
/03366029
/fr
/900164000042046-2
/900164000040674-2
/0000/page/4988
/0900164000209185
/99100100063779
/page/17586
/20003512547
/newstl/amitriptyline/lorazepam-vs-alprazolam-treatment-of-panic-disorder.php
/06739990
/985112000387339-2
/991001000333858/feed
/900164000037341-2
/page/11532
/900164001091015%C2%A0
/demo2/mike-arms-microchip-movement-2013/
/90016400106807
/03527581
/20001061816
/07063125
/73024317
/900164000043632-2
/054*315*769
/newstl/amitriptyline/xanax-4mg-bars.php
/900164000042053-3
/0000/page/23289
/0000/page/27088
/900164000043714-3
/20001715924
/90004000502244
/page/20172
/index.php
/AVID009-294-777
/005535885/feed
/900164000703268x%60
/15126857
/wp-content/sedlex/backup-scheduler/map.php
/04306026
/991001000244838/feed
/900164000987551c
/41326333
/900164000043025-2
/0000/page/29307
/900164001065796/feed
/newstl/amitriptyline/soma-350-mg-get-high.php
/991001000414409c
/page/18317
/page/23462
/0000/page/10983
`
)
