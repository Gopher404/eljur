function formatDate(month, day) {
    if (day === 100) {
        return "Итог"
    }
    let dayS = `${day}`
    if (dayS.length === 1) {
        dayS = "0" + dayS
    }
    let monthS = `${month}`
    if (monthS.length === 1) {
        monthS = "0" + monthS
    }
    return monthS + "." + dayS
}

function handleResponseCode(code, urlToLogin, message) {
    if (code !== 200) {
        alert(message)
        if (code === 401) {
            window.location.replace(urlToLogin)
        }
        return false
    }
    return true
}