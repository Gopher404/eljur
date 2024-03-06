
window.document.onkeydown = function (ev) {
    startRefocus(ev);
}
function startRefocus(event) {
    let  key = event.keyCode;
    let targetElement = event.target || event.srcElement;
    focusMe(targetElement, key);
}
function focusMe(input, key) {
    let needFocusElement = false;
    function detectColumn(td) {
        let result = 0, x;
        while (td = td.previousElementSibling) {
            result++
        }
        return result;
    }
    try {
        switch (key) {
            case 37:
                needFocusElement = input.parentNode.previousElementSibling.querySelector("input");
                break;
            case 39:
                needFocusElement = input.parentNode.nextElementSibling.querySelector("input");
                break;
            case 38:
                needFocusElement = input.parentNode.parentNode.previousElementSibling.querySelectorAll("td")[detectColumn(input.parentNode)].querySelector("input");
                break;
            case 40:
                needFocusElement = input.parentNode.parentNode.nextElementSibling.querySelectorAll("td")[detectColumn(input.parentNode)].querySelector("input");
                break
            default:
                needFocusElement = false
        }
    } catch (e) {
        needFocusElement = false;
    }
    if (!needFocusElement) return;
    needFocusElement.focus();
}