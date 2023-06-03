// gets the key from captcha
// and renders one more captcha

var captcha = document.querySelectorAll('iframe[title="reCAPTCHA"]')[0]
var url = new URL(captcha.src)
var key = url.searchParams.get("k")
console.log(key)

grecaptcha.ready(function () {
    var el = document.createElement("div")
    document.querySelector("body").appendChild(el)
    grecaptcha.render(el, {
        sitekey: key,
    })
})
