$(function() {
    //deal with main page form submission
    let $torrents_table_id = $('#table_torrents');
    let $torrents_table;
    let $inputNameMovies = $('#inputNameMovies');
    resetTorrentDataTable();
    //$('.dataTables_length').addClass('bs-select');


    initPopovers();


    //
    $("#form_search_movies").on("submit",function (e) {e.preventDefault();});
    $inputNameMovies.on('keypress',function(e) {
        if(e.which === 13) {
            $("#main_submit_movies").trigger('click');
        }
    });
    $inputNameMovies.on("input",function() {
        $("#inputResourceImdbID").val("");
    });
    $("#form_search_torrents").on("submit",function (e) {e.preventDefault();});
    $('#inputNameTorrents').on('keypress',function(e) {
        if(e.which === 13) {
            $("#main_submit_torrents").trigger('click');
        }
    });



    $("#main_submit_movies").on("click",function (e) {
        e.preventDefault();
        let inputName = $("#inputNameMovies");
        let inputResourceImdbID = $("#inputResourceImdbID");

        if (inputName.val().length === 0){
            notifyErrorOnDiv(inputName ,"Please specify a movie name","right");
        }else{
            let singleResultDiv = $('#div_single_result');
            customHide(singleResultDiv);
            displayLoadingSubmit();

            if (inputResourceImdbID.val().length > 0 && !isMobile){
                $.ajax({
                    url : "/central",
                    type : 'POST',
                    cache : false,
                    data : {
                        ajax : true,
                        resourceName : inputName.val(),
                        resourceImdbID : inputResourceImdbID.val(),
                        action : "getResourceInfo",
                    },
                    success : function (result) {
                        hideLoadingSubmit();
                        if (result.hasOwnProperty('error')){
                            notifyErrorOnDiv(".main-submit-button","resource unavailable")
                        }else {
                            if (result.hasOwnProperty('response')){
                                let resource = result['response'];

                                /*
                                * now we have our result. If it has more than 1 element, than we have to display it using 'div_torrents'.
                                * otherwise we can use 'div_single_result'
                                * */



                                let singleResultDivBanner = $('#div_result_banner');
                                let singleResultDivInfo = $('#div_result_info');
                                let singleResultDivTorrents = $('#div_result_torrents');

                                //  SET RESULT BANNER
                                if (resource.hasOwnProperty('Poster')) {
                                    $('#img_resource_banner').attr("src", resource["Poster"])
                                } else {
                                    console.error("single resource has no poster")
                                }


                                //  SET RESULT INFO
                                $('#single_result_title').text(resource['Title']);
                                $('#single_result_year').text(resource['Year']);
                                $('#single_result_categories').text(resource["Genre"]);
                                $('#single_result_stars').text(resource["Actors"]);
                                $('#single_result_directors').text(resource["Director"]);

                                customShow(singleResultDiv);

                                fetchTorrent(inputName.val(),inputResourceImdbID.val(),"movie");
                            } else {
                                //TODO handle error, check first for 'error'

                                notifyErrorOnDiv(".main-submit-button","resource unavailable")
                            }
                        }
                    },
                    error: dealWithAjaxError,
                });
            }else{
                fetchTorrent(inputName.val(),inputResourceImdbID.val(),"movie");
            }
        }
    });

    $("#main_submit_torrents").on("click",function (e) {
        e.preventDefault();
        let inputName = $("#inputNameTorrents");
        if (inputName.val().length === 0){
            notifyErrorOnDiv(inputName ,"Please specify a torrent name","right");
        }else{
            let singleResultDiv = $('#div_single_result');
            customHide(singleResultDiv);
            fetchTorrent(inputName.val());
        }
    });

    $inputNameMovies.autocomplete({
        source: function(request, response) {
            $.ajax({
                url : "/central",
                type : 'POST',
                cache: false,
                data : {
                    ajax : true,
                    resourceName : $("#inputNameMovies").val(),
                    action : "suggest",
                },
                success : function (result) {
                    if (result.hasOwnProperty("response")){
                        let queryResponse = result["response"];
                        let suggestions = [];
                        if (queryResponse instanceof Array){
                            /*$.each(queryResponse, function(suggestionIdx, suggestion){
                                if (suggestion.hasOwnProperty("resource_id") && suggestion.hasOwnProperty("suggestion_value")){
                                    let resourceId = suggestion["resource_id"];
                                    let suggestionValue = suggestion["suggestion_value"];
                                    suggestions.push({"label":suggestionValue,"value":resourceId});
                                }else{
                                    console.error("properties 'resource_id' and 'suggestion_value' not found in suggestion");
                                    return false;
                                }
                            });*/
                            $.each(queryResponse, function(suggestionIdx, suggestion){
                                if (suggestion.hasOwnProperty("imdbID") && suggestion.hasOwnProperty("Title")){
                                    let resourceId = suggestion["imdbID"];
                                    let suggestionValue = suggestion["Title"] + " ("+suggestion["Year"]+")";
                                    suggestions.push({"label":suggestionValue,"value":resourceId});
                                }else{
                                    console.error("properties 'imdbID' and 'Title' not found in suggestion");
                                    return false;
                                }
                            });
                        }
                        response(suggestions)
                    }else{
                        console.error("property 'response' not found in suggestion");
                        response([])
                    }
                }
            });
        },
        select: function (event, ui) {
            // Set selection
            $("#inputNameMovies").val(ui.item.label); // display the selected text value
            $("#inputResourceImdbID").val(ui.item.value); // display the selected text ID
            return false;
        },
        change: function(event, ui){
            $("#inputResourceImdbID").val(ui.item? ui.item.value : "");
        },
        minLength: 3
    });



    function fetchTorrent(iKeyword, iImdbID, iType){
        displayLoadingSubmit();
        let torrentsDiv = $('#div_torrents');
        customHide(torrentsDiv);
        $.ajax({
            url : "/central",
            type : 'POST',
            cache : false,
            data : {
                ajax : true,
                keyword : iKeyword,
                imdbID : iImdbID,
                type: iType,
                action : "getTorrents",
            },
            success : function (result) {
                hideLoadingSubmit();
                if (result.hasOwnProperty('response')){
                    let torrents = result['response'];

                    /*
                    * now we have our result. If it has more than 1 element, than we have to display it using 'div_torrents'.
                    * otherwise we can use 'div_single_result'
                    * */

                    resetTorrentDataTable();
                    if (torrents instanceof Array){
                        if (torrents.length > 0){
                            //now fetch and populate torrent table
                            for(let i=0 ; i<torrents.length ; i++){
                                let currentTorrent = torrents[i];
                                if (currentTorrent.hasOwnProperty("name") && currentTorrent.hasOwnProperty("size") && currentTorrent.hasOwnProperty("magnet_link")
                                    && currentTorrent.hasOwnProperty("peers") && currentTorrent.hasOwnProperty("files")){

                                    let newRow = $torrents_table.row.add([
                                        currentTorrent["name"],
                                        humanFileSize(currentTorrent["size"]),
                                        currentTorrent["peers"],
                                        '<a class="magnet-link" href="'+currentTorrent["magnet_link"]+'" ' +
                                        'rel="popover" ' +
                                        'data-trigger="hover" ' +
                                        'data-original-title="<a class=\'magnet-link-popup-header\'><strong>Magnet Link</strong></a>" ' +
                                        'data-content="' +
                                        '<a class=\'magnet-link-popup-body\'>' +
                                            'Clicking this link will open your default torrent BitTorrent client (eg. qBittorrent, Transmission, uTorrent etc...) to start ' +
                                            'the download.' +
                                            '<br /> Do you still not have a torrent client? Check these out!' +
                                            '<br /> <a href=\'https://www.qbittorrent.org\' target=\'_blank\'><b>qBittorrent<b/></a>' +
                                            '<br /> <a href=\'https://transmissionbt.com/download\' target=\'_blank\'><b>Transmission<b/></a>' +
                                            '<br /> <a href=\'https://www.utorrent.com\' target=\'_blank\'><b>uTorrent<b/></a>' +
                                        '</a> "' +
                                        ' data-html="true"></a>',
                                        formatFiles(currentTorrent["files"])
                                    ]);

                                    //no poster on mobile, if's annoying
                                    if (!isMobile){
                                        if (currentTorrent.hasOwnProperty("poster")){
                                            let poster = currentTorrent["poster"];
                                            if (poster.length > 0){
                                                $(newRow).attr('data-toggle','popover-poster');
                                                $(newRow).attr('data-img',poster);
                                                //pointer-events: none
                                            }
                                        }
                                    }
                                }else{
                                    console.log("currentTorrent is missing some property",currentTorrent);
                                }
                            }
                            //finally, draw the main torrents table
                            $torrents_table.draw();

                            initPopovers();

                            customShow(torrentsDiv);

                            //bring user to the result
                            $('html, body').animate({
                                scrollTop: $("#div_results").offset().top
                            }, 500);
                        }else{
                            notifyErrorOnDiv(".main-submit-button" ,"resource not available");
                        }
                    }else{
                        //it always has to be an array. If it is not, this is an error
                        //TODO handle error
                        notifyErrorOnDiv($(".main-submit-button") ,"no torrent available","right");
                    }
                }else{
                    hideLoadingSubmit();
                    //TODO handle error, check first for 'error'

                    notifyError("some error occurred")
                }
            },
            error: dealWithAjaxError,
        });
    }
    function formatFiles (inputFilesList) {
        if (inputFilesList instanceof Array){
            //creating table listing all files and respective size
            let table =
                '<table class="table table-hover table-striped" cellspacing="0" width="100%">' +
                '            <thead>' +
                '            <tr>' +
                '                <th>Path</th>' +
                '                <th>Size</th>' +
                '            </tr>' +
                '            <tbody>'
            ;

            for(let i = 0 ; i < inputFilesList.length ; i++){
                table +=
                    '<tr>'+
                        '<td>'+inputFilesList[i].path+'</td>'+
                        '<td>'+humanFileSize(inputFilesList[i].size)+'</td>'+
                    '</tr>'
            }

            table +=
                '            </tbody>'+
                '            </thead>' +
                '</table>';

            return table;
        }else{
            return ""
        }
    }

    function initPopovers() {
        initStandardPopover();
        initPosterPopover();
    }
    function initStandardPopover(){
        $('[rel="popover"]').popover(
            {
                html: true,
                trigger: 'manual',
            })
            .on('mouseenter', function () {
                let _this = this;
                $(this).popover('show');
                $('.popover').on('mouseleave', function () {
                    $(_this).popover('hide');
                });
            })
            .on('mouseleave', function () {
                let _this = this;
                setTimeout(function () {
                    if (!$('.popover:hover').length) {
                        $(_this).popover('hide');
                    }
                }, 300);
            });
    }
    function initPosterPopover(){
        $('[data-toggle="popover-poster"]').popover({
            html: true,
            trigger: 'hover',
            content: function () {
                return "<img class='hover-img' src='" + $(this).data('img') + "' alt='hover--poster-img'/>";
            },
            template:getPopoverCustomTemplate("popover-poster")
        })
    }

    function resetTorrentDataTable(){
        $torrents_table_id.DataTable().clear();
        $torrents_table_id.DataTable().destroy();
        $torrents_table = $torrents_table_id.DataTable(
            {
                "aaSorting": [], //not sorting initially, preserving DB order (the user can choose after)
                responsive: true,
                "columnDefs":
                [
                    { "width": "40%", "targets": 0 },
                    { "width": "20%", "targets": [1,2,3] },

                    { className: "tName", "targets": 0 },
                    { className: "tSize", "targets": 1 },
                    { className: "tPeers", "targets": 2 },
                    { className: "tDownload", "targets": 3 }
                ],
                drawCallback: function () {
                    initPosterPopover()
                }
            }
        ).columns.adjust().responsive.recalc();
    }

    function getPopoverCustomTemplate(className) {
        return '<div class="popover '+className+'" role="tooltip"><div class="arrow"></div><h3 class="popover-header"></h3><div class="popover-body"></div></div>';
    }

});
