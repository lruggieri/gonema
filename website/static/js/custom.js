$( document ).ready(function() {
    //deal with main page form submission
    $("#submit_main_form").click(function (e) {
        if ($("#inputResourceName").val().length === 0 && $("#inputResourceImdbID").val().length === 0){
            e.preventDefault();
            $("#submit_main_form").notify(
                "Please specify at least 1 value",
                {
                    position:"bottom center",
                    showAnimation: 'slideDown',
                    showDuration: 400,
                    autoHide: true,
                    // if autoHide, hide after milliseconds
                    autoHideDelay: 2000,
                }
            );
        }else{
            e.preventDefault();
            $.ajax({
                url : "/central",
                type : 'POST',
                cache : false,
                data : {
                    ajax : true,
                    resourceName : encodeURIComponent($('#inputResourceName').val()),
                    resourceImdbID : encodeURIComponent($('#inputResourceImdbID').val()),
                    action : "getResourceInfo",
                },
                success : function (result) {
                    /*var $resultDiv = $('#div_results');
                    $resultDiv.empty();

                    var $imageCol = $('<div class="col-md-4"></div>');
                    var $slideShowCol = $('<div class="col-md-8"></div>');

                    $resultDiv.append($imageCol);
                    $resultDiv.append($slideShowCol);*/

                    if (result.hasOwnProperty('resources')){
                        let resources = result['resources'];

                        /*
                        * now we have our result. If it has more than 1 element, than we have to display it using 'div_slideshow_results'.
                        * otherwise we can use 'div_single_result'
                        * */

                        let singleResultDiv = $('#div_single_result');
                        let slideShowResultsDiv = $('#div_slideshow_results');
                        customHide(singleResultDiv);
                        customHide(slideShowResultsDiv);

                        if (resources instanceof Array){
                            if (resources.length > 1){
                                //use 'div_slideshow_results'
                                //TODO
                            }else{
                                //use 'div_single_result'

                                let singleResultDivBanner = $('#div_result_banner');
                                let singleResultDivInfo = $('#div_result_info');
                                let singleResultDivTorrents = $('#div_result_torrents');
                                let singleResult = resources[0];


                                //  SET RESULT BANNER
                                if (singleResult.hasOwnProperty('images')){
                                    let images = singleResult['images'];
                                    if (images.hasOwnProperty('big')){
                                        $('#img_resource_banner').attr("src",images["big"])
                                    }else{
                                        console.error("single resource has no image set")
                                    }
                                }else{
                                    console.error("single resource has no images")
                                }

                                //  SET RESULT INFO
                                $('#single_result_title').text(singleResult['title']);
                                $('#single_result_year').text(singleResult['year']);
                                if ( singleResult['categories'] instanceof Array){
                                    $('#single_result_categories').text(singleResult['categories'].filter(Boolean).join(", "));
                                }
                                if ( singleResult['stars'] instanceof Array){
                                    $('#single_result_stars').text(singleResult['stars'].filter(Boolean).join(", "));
                                }
                                if ( singleResult['directors'] instanceof Array){
                                    $('#single_result_directors').text(singleResult['directors'].filter(Boolean).join(", "));
                                }


                                // SET RESULT TORRENT SLIDESHOW
                                if (singleResult['available_torrents'] instanceof Array){
                                    if (singleResultDivTorrents.hasClass('slick-initialized')){
                                        singleResultDivTorrents.slick('unslick');
                                        singleResultDivTorrents.empty(); //reset torrent slideshow
                                    }
                                    singleResultDivTorrents.slick({
                                        infinite: false,
                                        speed: 200,
                                        slidesToShow: 3,
                                        slidesToScroll: 1,
                                    });


                                    $.each(singleResult['available_torrents'], function(key, value){
                                        let newTorrentDiv = $('<div></div>');

                                        let infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Name:</h4>'));
                                        infoItemDiv.append($('<a">'+value["name"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Quality:</h4>'));
                                        infoItemDiv.append($('<a">'+value["quality"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Length:</h4>'));
                                        infoItemDiv.append($('<a">'+value["length"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Resolution:</h4>'));
                                        infoItemDiv.append($('<a">'+value["resolution"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Size:</h4>'));
                                        infoItemDiv.append($('<a">'+value["size"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Sound:</h4>'));
                                        infoItemDiv.append($('<a">'+value["sound"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Codec:</h4>'));
                                        infoItemDiv.append($('<a">'+value["codec"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Seeders:</h4>'));
                                        infoItemDiv.append($('<a">'+value["seeders"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Leechers:</h4>'));
                                        infoItemDiv.append($('<a">'+value["leechers"]+'</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        infoItemDiv = $('<div class="div_info_item"></div>');
                                        infoItemDiv.append($('<h4 class="inline">Link:</h4>'));
                                        infoItemDiv.append($('<a href="'+value["magnet_link"]+'">torrent</a>'));
                                        newTorrentDiv.append(infoItemDiv);

                                        singleResultDivTorrents.slick('slickAdd',newTorrentDiv);
                                    });
                                }

                                customShow(singleResultDiv);
                            }
                        }else{
                            //it always has to be an array. If it is not, this is an error
                            //TODO handle error
                            console.error("resources is not an array")
                        }

                    }else{
                        //TODO handle error, check first for 'error'
                        console.error('error to be handled')
                    }

                }
            });


            $("#submit_main_form").notify(
                "Not yet implemented, wait for it...",
                {
                    style: 'simplegreen',
                    position:"bottom center",
                    showAnimation: 'slideDown',
                    showDuration: 400,
                    autoHide: true,
                    // if autoHide, hide after milliseconds
                    autoHideDelay: 2000,
                }
            );
        }
    });


    function customShow($inputDiv){
        $inputDiv.show();
    }
    function customHide($inputDiv){
        $inputDiv.hide();
    }

    $.notify.addStyle('simplegreen', {
        html: "<div><span data-notify-text/></div>",
        classes: {
            base: {
                "white-space": "nowrap",
                "padding": "5px"
            },
            superblue: {
                "color": "green",
            }
        }
    });
});