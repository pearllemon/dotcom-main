window.mainstl=function($){

	var SelectorFormPoster = "form.poster";
	var SelectorFormSubmitter=".poster-submitter";
	var selectorErrorContainer=".js-error-container";
	var selectorErrorText=".js-error-text";
	var selectorSuccessContainer=".js-success-container";
	var selectorSuccessText=".js-success-text";
	var selectorErrorResponse=".js-error-response";
	var selectorLoading=".js-loading";
	var selectorAsyncContent=".js-async-content";
	var selectorFormPaypal = "form#paypal-form";
	var selectorSubmitterPaypal =".submitter-paypal";
	var dataTemplate="template";
	var dataValue="value";
	var dataContainer="container";
	var dataAsyncUrl="asyncUrl";

	var dataContinueUrl = "continueurl";
	var dataFormToSubmit = "form"; 
	var dataScrollTop="scrolltop";

	function Send(form,url,continueURL,progressPopup,scrollTop) {
		var options = {};
		var enctype= form.attr("enctype");
		var data;
		var isMultipart=false;
		if (enctype==="multipart/form-data"){
			data=new FormData( form.get(0) );
			isMultipart=true;
		}else{
			data=form.serialize();
		}
		
		// reset error and success panels
		$(selectorErrorContainer).hide();
		$(selectorSuccessContainer).hide();

		if (progressPopup !== "") {
			$(progressPopup).show();
		}
	    if (continueURL==="") {
	    	options = {
	            "success": function(response) {
	            	if (progressPopup !== "") {
	            		$(progressPopup).hide();
	            	}

	            	console.log(response)
	            	if (response=="reload") {
	            		window.location.reload(true);
	            		return
	            	}
	            	if (response=="paypal") {
	            		$(selectorFormPaypal).submit();
	            		return
	            	}
	            	if(response.indexOf("/") == 0) {
	            		window.location.href=response;
	            		return
	            	}

	            	$(selectorSuccessText).html(response);
	            	$(selectorSuccessContainer).show();
	            	form.get(0).reset();
	            	
	            	
	            },
	            "error":   function(jqXHR, textStatus, errorThrown) {
	            	if (progressPopup !== "") {
	            		$(progressPopup).hide();
	            	}
	            	var response = $($.parseHTML(jqXHR.responseText));
	            	$(selectorErrorText).html(response.filter(selectorErrorResponse));
	            	$(selectorErrorContainer).show();
            	},
	            "method":  "POST",
	            "data":    data,
	            "url": url
	        };
	    }else{
	    	options = {
	            "success": function() {
	            	if (progressPopup !== "") {
	            		$(progressPopup).hide();
	            	}
	            	location.href = continueURL;
	            	$(selectorLoading).hide();
	            },
	            "error":   function(jqXHR, textStatus, errorThrown) {
	            	if (progressPopup !== "") {
	            		$(progressPopup).hide();
	            	}
	            	var response = $($.parseHTML(jqXHR.responseText));
	            	$(selectorErrorText).html(response.filter(selectorErrorResponse));
	            	$(selectorErrorContainer).show();

            	},
	            "method":  "POST",
	            "data":    data,
	            "url": url
	        };
	    }
	    if (isMultipart) {
	    	options["contentType"]=false;
	    	options["processData"]=false;
	    }

	    if (scrollTop=="true") {
	    	$('html, body').animate({ scrollTop: 200 }, 'fast');
	    }
		$(selectorErrorContainer).hide();
        $.ajax(options);
        
	}

	function onFormSubmit (evt) {
		evt.preventDefault();
		var form =  $(evt.target);
		var continueUrl = form.data(dataContinueUrl);
		var scrollTop=form.data(dataScrollTop);
		var url = form.attr("action");
		if(continueUrl) {
			Send(form,url, continueUrl, selectorLoading,scrollTop);
		}else{
			Send(form, url,"", selectorLoading,scrollTop);
		}
	}

	function onClickedSubmit (evt) {
		evt.preventDefault();
		var clickedElt =  $(evt.currentTarget);
		var formSelector = clickedElt.data(dataFormToSubmit);
		if (formSelector !== "") {
			console.log("submit form "+formSelector);
			$("form"+formSelector).submit();
		}else{
			logger.warn("data formSelector is empty!");
		}
	}

	// htmlAsyncAjaxGet remplace le contenu de l'element html "obj" par un contenu html
	// résultat d'une requete AJAX
	function htmlAsyncAjaxGet(obj) {
		var url= obj.data(dataAsyncUrl);
		options = {
	         "success": function(html) {
	            obj.html(html);
	         },
	         "error":   function(jqXHR, textStatus, errorThrown) {
		    	obj.html("ERR");
	         },
	         "method":  "GET",
	         "url": url
	    };
	    $.ajax(options);
	}

	function clean_text (text){
        text = text.replace(/"/g, '""');
        return '"'+text+'"';
    };

    function doSearch(evt,force) {
	    	
    	var search = $("input#search").val();
    	if (!force && search.length < 3) return;
    	jQuery(selectorLoading).show();
	    var request = {
	      type: "GET",
	      url: "https://www.savethislife.com/searchmicrochip?search="+search,
	      success: function (data,textStatus,jqXHR) {
	            console.log("success");
	            jQuery(selectorLoading).hide();
	            jQuery("#search-results").html(jqXHR.responseText);
	                
	        } ,
	      error: function (jqXHR, textStatus, errorThrown) {
				console.log("ajax error:" + errorThrown);
				jQuery(selectorLoading).hide();
	      		console.log(jqXHR.responseText);
				jQuery("#search-results").html("<font color='red'>Error: Cannot search microchip. Please contact us to fix this. Thanks!</font>");
			},
	    };
	    jQuery.ajax(request);
    }



	function init() { 
	    $("body").on("submit", SelectorFormPoster, onFormSubmit );
        $("body").on("click", SelectorFormSubmitter, onClickedSubmit );

        // asyncload
        $(selectorAsyncContent).each(function(index,value){
		  	htmlAsyncAjaxGet($(this));
		})

		$("body").on("click",".js-checkout",function (evt) {
			$("#checkout-div").show();
			$("form#shopping_form").hide();
			var clickedElt =  $(evt.currentTarget);
			clickedElt.hide();
			$('html, body').animate({ scrollTop: 200 }, 'fast');
		});

		$("body").on('click',selectorSubmitterPaypal,function (evt) {
			evt.preventDefault();
			$("form#checkout_form").submit();
			return false;
		});

		$("body").on('click',".search-button",doSearch);
	    $("body").on('keyup',"input#search",function(e) {
			clearTimeout($.data(this, 'timer'));
			if (e.keyCode == 13)
			  doSearch(e,true);
			else
			  $(this).data('timer', setTimeout(doSearch, 500));
		});
	}

	$(function(){ 
	    console.log("start main stl");
	    init();
	});

	var mainstl= {};

	return mainstl;

}(jQuery);