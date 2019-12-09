function dealWithAjaxError(request, status, error) {
    hideLoadingSubmit();
    notifyError("Service unavailable. Sorry for the inconvenience.");
    console.log("StatusCode "+request.status+", response:"+request.responseText);
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