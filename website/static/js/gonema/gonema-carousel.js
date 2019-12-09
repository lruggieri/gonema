$(function() {
    let maxSlidePerCarousel = 5;
    let carouselClass =  $('.carousel');
    //carouselClass.carousel(); //unnecessary because of  data-ride
    carouselClass.swipe({
        swipe: function (event, direction, distance, duration, fingerCount, fingerData) {
            if (direction === 'left') $(this).carousel('next');
            if (direction === 'right') $(this).carousel('prev');
        },
        allowPageScroll: "vertical"
    });

    $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
        let target = $(e.target).attr("id"); // activated tab
        if (target === 'torrents-tab'){
            customHide($('#div_slideshow'))
        }else{
            customShow($('#div_slideshow'))
        }
    });

    /*$('div.carousel-inner div.carousel-item div.row div.col-sm img').on('click',function(){
        console.log("click detected");
        //console.log($(this).data('imdbID'))
    });
    */
    /*$('.carousel-poster').on('click',function(){
        console.log("click detected on img");
        //console.log($(this).data('imdbID'))
    });
    $('.carousel-item').on('click',function(){
        console.log("click detected on item");
        //console.log($(this).data('imdbID'))
    });*/
    $('a[href="#idparticolare"]').on('click',function(e){
        //e.preventDefault();
        console.log("click detected on href!");
        //console.log($(this).data('imdbID'))
    });
    /*$('.poster-link').on('click',function(e){
        e.preventDefault();
        console.log("click detected on inner");
        //console.log($(this).data('imdbID'))
    });*/
    $(document).on("click", ".poster-link", function(e){
        e.preventDefault();
        let posterName = $(e.target).data('name');
        let posterImdbID = $(e.target).data('imdbid');

        $("#inputNameMovies").val(posterName); // display the selected text value
        $("#inputResourceImdbID").val(posterImdbID); // display the selected text ID
        $("#main_submit_movies").trigger('click');
    });




    //ajax to fill each carousel

    $.ajax({
        url : "/central",
        type : 'POST',
        cache : false,
        data : {
            ajax : true,
            aggType : "most_present",
            resType : "movie",
            action : "getAggregations",
        },
        success : function (result) {
            buildCarousel('#carouselMoviesMostPresent',result)
        },
        error: dealWithAjaxError,
    });
    $.ajax({
        url : "/central",
        type : 'POST',
        cache : false,
        data : {
            ajax : true,
            aggType : "most_present",
            resType : "serie",
            action : "getAggregations",
        },
        success : function (result) {
            buildCarousel('#carouselSeriesMostPresent',result)
        },
        error: dealWithAjaxError,
    });
    $.ajax({
        url : "/central",
        type : 'POST',
        cache : false,
        data : {
            ajax : true,
            aggType : "most_shared",
            resType : "movie",
            action : "getAggregations",
        },
        success : function (result) {
            buildCarousel('#carouselMoviesMostShared',result)
        },
        error: dealWithAjaxError,
    });
    $.ajax({
        url : "/central",
        type : 'POST',
        cache : false,
        data : {
            ajax : true,
            aggType : "most_shared",
            resType : "serie",
            action : "getAggregations",
        },
        success : function (result) {
            buildCarousel('#carouselSeriesMostShared',result)
        },
        error: dealWithAjaxError,
    });

    function buildCarousel(iCarouselDivID, iResponse){
        if (iResponse.hasOwnProperty('error')){
            notifyErrorOnDiv(".main-submit-button","resource unavailable")
        }else {
            if (iResponse.hasOwnProperty('response')){
                let resp = iResponse['response'];

                if (resp instanceof Array){
                    let currentCarouselItem = $('<div class="carousel-item"></div>');
                    let currentRow = $('<div class="row"></div>');
                    let indicators = 0;
                    let insert = function () {
                            /*let href = $('<a href="http://estiloasertivo.blogspot.com.es/"> </a>');
                            href.appendTo(currentCarouselItem);
                            currentRow.appendTo(href);*/
                            currentRow.appendTo(currentCarouselItem);
                            currentCarouselItem.appendTo($(iCarouselDivID+' > .carousel-inner'));
                            //Indicators
                            $('<li data-target="'+iCarouselDivID+'" data-slide-to="'+indicators+'"></li>').appendTo($(iCarouselDivID+' > .carousel-indicators'));
                            indicators++;

                            //reset
                            currentCarouselItem = $('<div class="carousel-item"></div>');
                            currentRow = $('<div class="row"></div>');
                    };
                    for (let i = 0; i < resp.length; i++){
                        if (resp[i].hasOwnProperty('poster') && resp[i].hasOwnProperty('imdbID') && resp[i].hasOwnProperty('name')){
                            let col = $('<div class="col-sm"></div>');
                            let href = $('<a class="poster-link" href="#"</a>');
                            let poster = $('<img class="d-block w-100 carousel-poster" src="'+resp[i].poster+'" alt="'+i+' slide" ' +
                                'data-imdbid="'+resp[i].imdbID+'" data-name="'+resp[i].name+'">');


                            poster.appendTo(href);
                            href.appendTo(col);
                            col.appendTo(currentRow);

                            if (i !== 0 && (i+1) % maxSlidePerCarousel === 0){
                                insert();
                            }
                        }

                    }
                    if (currentRow.children().length > 0){
                        insert();
                    }


                    $(iCarouselDivID+' > .carousel-indicators > li').first().addClass('active');
                    $(iCarouselDivID+' > .carousel-inner > .carousel-item').first().addClass('active');
                    $(iCarouselDivID).bind('mousewheel', function(e) {
                        e.preventDefault();
                        if(e.originalEvent.wheelDelta /120 > 0) {
                            $(this).carousel('next');
                        } else {
                            $(this).carousel('prev');
                        }
                    });
                }
            } else {
                //TODO handle error, check first for 'error'

                notifyError("carousel not available at the moment")
            }
        }
    }

});

