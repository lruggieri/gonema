$(function() {

    //tracks ENTER pressed on the main search input box
    $('.main-search-box').on('keypress',function(e) {
        if(e.which === 13) { //on enter
            let searchValue =  $(this).val();
            if (searchValue.length > 0){
                window.dataLayer = window.dataLayer || [];
                window.dataLayer.push({
                    event: 'mainSearchEnter',
                    mainSearchInput: searchValue,
                });
            }
        }
    });

    //tracks a click on the torrent download button
    //datatable makes it difficult to track directly .class.on('click') ...
    $(document).on('click','.magnet-link',function(){
        window.dataLayer = window.dataLayer || [];
        window.dataLayer.push({
            event: 'magnetLinkClick',
        });
    });


    //tracks a click on a poster
    $(document).on('click','.carousel-poster',function(){
        let posterName = $(this).data('name');
        window.dataLayer = window.dataLayer || [];
        window.dataLayer.push({
            event: 'posterClick',
            posterName:posterName,
        });
    });

    // tracks a CLICK event on the main search form
    // triggered by: mainSearchEnter, posterClick
    $('.main-submit-button').on('click',function(){
        let searchValue = $(this).parent().parent().parent().find(".main-search-box").val();
        if(searchValue !== undefined && searchValue.length > 0){
            window.dataLayer = window.dataLayer || [];
            window.dataLayer.push({
                event: 'mainSearchClick',
                mainSearchInput: searchValue,
            });
        }
    });

});