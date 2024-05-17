

let domain = 'luoye-g.top'
/**
 * 登录使用
 */
function sign_in(event) {
    event.preventDefault();
    let email = document.getElementById('floatingInput')
    let password = document.getElementById('floatingPassword')

    var xhr = new XMLHttpRequest();
    xhr.open("POST", 'https://' + domain + '/api/login', false);
    xhr.setRequestHeader('Accept', 'application/json');

    var formData = new FormData();

    formData.append("email", email.value);
    formData.append("password", password.value);

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status == 200) {
                // 成功后直接进行页面跳转
                console.log('跳转');
                window.location.href = 'https://' + domain + '/test.html';
            } else {
                console.error('Server responded with status:', xhr.status);
            }
        }
    };

    xhr.onerror = function () {
        console.error('Request failed', xhr.statusText);
    };

    xhr.send(formData);
}