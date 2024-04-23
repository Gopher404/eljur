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
    return dayS + "." + monthS
}

function handleResponseCode(code, message) {
    if (code !== 200) {
        alert(message)
        if (code === 401) {
            window.location.replace("/login")
        }
        return false
    }
    return true
}

function LogOut() {
    deleteCookie("token");
    window.location.href = "/login"
}
