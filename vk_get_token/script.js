let lsKeys = Object.keys(localStorage);
let token = ""
for (let key of lsKeys) {
    if (key.indexOf("web_token") > -1) {
        token = JSON.parse(localStorage[key])["access_token"];
    }
}
return token