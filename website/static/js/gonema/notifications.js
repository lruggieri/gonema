function notifyError(iMessage, iCustomHeader) {
    let header = "Error";
    if (iCustomHeader !== undefined && iCustomHeader.length > 0) {
        header = iCustomHeader;
    }
    $.notify({
        title: '<strong>' + header + '</strong>',
        message: iMessage,
    }, {
        type: 'danger'
    });
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