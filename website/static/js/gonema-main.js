$(function() {
    //deal with main page form submission
    let $torrents_table_id = $('#table_torrents');
    let $torrents_table;
    resetTorrentDataTable();
    $('.dataTables_length').addClass('bs-select');


    $(function () {
        $('[rel="popover"]').popover(
            {
                gpuAcceleration: !(window.devicePixelRatio < 1.5 && /Win/.test(navigator.platform))
            }
        )
    });



    $("#form_search_movies").on("submit",function (e) {
        e.preventDefault();
        $("#main_submit_movies").trigger('click');
    });
    $( "#inputNameMovies" ).on("input",function() {
        $("#inputResourceImdbID").val("");
    });
    $("#form_search_torrents").on("submit",function (e) {
        e.preventDefault();
        $("#main_submit_torrents").trigger('click');
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

            if (inputResourceImdbID.val().length > 0){
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

                                fetchTorrent(inputName.val(),"movie");
                            } else {
                                //TODO handle error, check first for 'error'

                                notifyErrorOnDiv(".main-submit-button","resource unavailable")
                            }
                        }
                    },
                    error: dealWithAjaxError,
                });
            }else{
                fetchTorrent(inputName.val(),"movie");
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
        change: function(event, ui){
            $("#inputResourceImdbID").val(ui.item? ui.item.value : "");
        },
        minLength: 3
    });



    function fetchTorrent(iKeyword, iType){
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
                                let newRow = $torrents_table.row.add([
                                    currentTorrent["name"],
                                    humanFileSize(currentTorrent["size"]),
                                    '<a class="magnet-link" href="'+currentTorrent["magnet_link"]+'"></a>',
                                    currentTorrent["peers"],
                                    formatFiles(currentTorrent["files"])
                                ]).draw().node();
                                if (currentTorrent.hasOwnProperty("poster")){
                                    let poster = currentTorrent["poster"];
                                    if (poster.length > 0){
                                        $(newRow).attr('data-toggle','popover-hover');
                                        $(newRow).attr('data-img',poster);
                                    }
                                }

                            }
                            $('[data-toggle="popover-hover"]').popover({
                                html: true,
                                trigger: 'hover',
                                content: function () {
                                    console.log($(this).data('img'));
                                    return '<img src="' + $(this).data('img') + '" />';
                                }
                            });
                            customShow(torrentsDiv);
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
    function notifyErrorOnDiv($notificationDiv, $message, $position){

        if (typeof $position === "undefined" || !$position.length){
            $position = "bottom center"
        }

        $($notificationDiv).notify(
                $message,
                {
                    position:$position,
                    showAnimation: 'slideDown',
                    showDuration: 400,
                    autoHide: true,
                    // if autoHide, hide after milliseconds
                    autoHideDelay: 2000,
                }
            );
    }

    function displayLoadingSubmit(iSubmitButton){
        $('.main-submit-button').attr("disabled",true);
        $('.spinner-submit').addClass("spinner-border spinner-border-sm")
    }
    function hideLoadingSubmit(iSubmitButton){
        $('.main-submit-button').removeAttr("disabled");
        $('.spinner-submit').removeClass("spinner-border spinner-border-sm")
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
                responsive: true,
                drawCallback: function() {
                    $('[data-toggle="popover-hover"]').popover({
                        html: true,
                        trigger: 'hover',
                        content: function () {
                            return "<img class='hover-img' src='" + $(this).data('img') + "'/>";
                        }
                    });
                }
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
        hideLoadingSubmit();
        notifyError("Service unavailable. Sorry for the inconvenience.");
        console.log("StatusCode "+request.status+", response:"+request.responseText);
    }
});
