$(function() {
    //deal with main page form submission
    let $torrents_table_id = $('#table_torrents');
    let $torrents_table;
    resetTorrentDataTable();

    $('.dataTables_length').addClass('bs-select');

    $("#main_submit_movies").on("click",function (e) {
        let inputName = $("#inputNameMovies");
        let inputResourceImdbID = $("#inputResourceImdbID");
        if (inputResourceImdbID.val().length === 0){
            e.preventDefault();
            notifyError(inputResourceImdbID ,"Please specify at least an Imdb ID","right");
        }else{
            e.preventDefault();
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
                    if (result.hasOwnProperty('error')){
                        notifyError(".main-submit-button","resource unavailable")
                    }else {
                        if (result.hasOwnProperty('response')){
                            let resource = result['response'];

                            /*
                            * now we have our result. If it has more than 1 element, than we have to display it using 'div_torrents'.
                            * otherwise we can use 'div_single_result'
                            * */

                            let singleResultDiv = $('#div_single_result');
                            customHide(singleResultDiv);

                            /*if (resources instanceof Array){
                                if (resources.length > 1){
                                    //use 'div_torrents'
                                    //TODO
                                    console.log("got more than 1 resource from a single ImdbID")
                                }else if (resources.length > 0){
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
                                    if ( singleResult['genre'] instanceof Array){
                                        $('#single_result_categories').text(singleResult['genre'].filter(Boolean).join(", "));
                                    }
                                    if ( singleResult['actors'] instanceof Array){
                                        $('#single_result_stars').text(singleResult['actors'].filter(Boolean).join(", "));
                                    }
                                    if ( singleResult['directors'] instanceof Array){
                                        $('#single_result_directors').text(singleResult['directors'].filter(Boolean).join(", "));
                                    }

                                    customShow(singleResultDiv);

                                    fetchTorrent(inputName.val());
                                }else{
                                    notifyError(".main-submit-button" ,"resource not available");
                                }
                            }else{
                                //it always has to be an array. If it is not, this is an error
                                //TODO handle error
                                notifyError(inputName ,"some error occurred","right");
                                console.error("resources is not an array")
                            }*/

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

                            fetchTorrent(inputName.val());
                        } else {
                            //TODO handle error, check first for 'error'

                            notifyError(".main-submit-button","resource unavailable")
                        }
                    }
                },
                error: dealWithAjaxError,
            });
        }
    });

    $("#main_submit_torrents").on("click",function (e) {
        e.preventDefault();
        let inputName = $("#inputNameTorrents");
        if (inputName.val().length === 0){
            e.preventDefault();
            notifyError(inputName ,"Please specify a torrent name","right");
        }else{
            let singleResultDiv = $('#div_single_result');
            customHide(singleResultDiv);
            fetchTorrent(inputName.val());
        }
    });

    function fetchTorrent(iKeyword){
        $.ajax({
            url : "/central",
            type : 'POST',
            cache : false,
            data : {
                ajax : true,
                keyword : iKeyword,
                action : "getTorrents",
            },
            success : function (result) {
                if (result.hasOwnProperty('response')){
                    let torrents = result['response'];

                    /*
                    * now we have our result. If it has more than 1 element, than we have to display it using 'div_torrents'.
                    * otherwise we can use 'div_single_result'
                    * */

                    let torrentsDiv = $('#div_torrents');
                    customHide(torrentsDiv);
                    resetTorrentDataTable();

                    if (torrents instanceof Array){
                        if (torrents.length > 0){

                            //now fetch and populate torrent table
                            for(let i=0 ; i<torrents.length ; i++){
                                let currentTorrent = torrents[i];
                                $torrents_table.row.add([
                                    currentTorrent["name"],
                                    humanFileSize(currentTorrent["size"]),
                                    '<a href="'+currentTorrent["magnet_link"]+'">Link</a>',
                                    currentTorrent["peers"],
                                    formatFiles(currentTorrent["files"])
                                ]).draw();
                            }
                            customShow(torrentsDiv);

                        }else{
                            notifyError(".main-submit-button" ,"resource not available");
                        }
                    }else{
                        //it always has to be an array. If it is not, this is an error
                        //TODO handle error
                        notifyError($(".main-submit-button") ,"no torrent available","right");
                    }

                }else{
                    //TODO handle error, check first for 'error'

                    notifyError(".main-submit-button","some error occurred")
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

    $("#inputNameMovies").autocomplete({
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
        minLength: 3
    });

    function notifyError($notificationDiv, $message, $position){

        if (typeof $position === "undefined" || !$position.length){
            $position = "bottom center"
        }

        $($notificationDiv).notify(
                $message,
                {
                    style: 'simplegreen',
                    position:$position,
                    showAnimation: 'slideDown',
                    showDuration: 400,
                    autoHide: true,
                    // if autoHide, hide after milliseconds
                    autoHideDelay: 2000,
                }
            );
    }

    function customShow($inputDiv){
        $inputDiv.show();
    }
    function customHide($inputDiv){
        $inputDiv.hide();
    }

    function resetTorrentDataTable(){
        $torrents_table_id.DataTable().clear();
        $torrents_table_id.DataTable().destroy();
        $torrents_table = $torrents_table_id.DataTable(
            {
                "aaSorting": [], //not sorting initially, preserving DB order (the user can choose after)
                responsive: true
            }
        );
    }
    function humanFileSize(bytes, si) {
        var thresh = si ? 1000 : 1024;
        if(Math.abs(bytes) < thresh) {
            return bytes + ' B';
        }
        var units = si
            ? ['kB','MB','GB','TB','PB','EB','ZB','YB']
            : ['KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB'];
        var u = -1;
        do {
            bytes /= thresh;
            ++u;
        } while(Math.abs(bytes) >= thresh && u < units.length - 1);
        return bytes.toFixed(1)+' '+units[u];
    }

    function dealWithAjaxError(request, status, error) {
        notifyError(".main-submit-button","Service unavailable. Sorry for the inconvenience.");
        console.log("StatusCode "+request.status+", response:"+request.responseText);
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
