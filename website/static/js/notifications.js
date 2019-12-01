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