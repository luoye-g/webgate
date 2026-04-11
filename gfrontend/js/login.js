// 自动检测输入类型函数
function detectInputType(identifier) {
    // 手机号正则：11位数字，1开头
    const phoneRegex = /^1[3-9]\d{9}$/;
    // 邮箱正则
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    
    if (phoneRegex.test(identifier)) {
        return 'phone';
    } else if (emailRegex.test(identifier)) {
        return 'email';
    } else {
        return 'unknown';
    }
}

// 登录函数
async function loginUser(event) {
    event.preventDefault();
    
    const loginMessage = document.getElementById('loginMessage');
    const identifier = document.getElementById('login-identifier').value;
    const password = document.getElementById('login-password').value;
    
    // 基础验证
    if (!identifier || !password) {
        showMessage('请输入手机号码/邮箱和密码', 'danger');
        return;
    }
    
    // 自动检测输入类型
    const inputType = detectInputType(identifier);
    
    if (inputType === 'unknown') {
        showMessage('请输入有效的手机号码或邮箱地址', 'danger');
        return;
    }

    // 显示加载状态
    const loginBtn = event.target;
    const originalText = loginBtn.innerHTML;
    loginBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 登录中...';
    loginBtn.disabled = true;

    try {
        // 调用登录API
        const response = await fetch('/api/user/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                identifier: identifier,
                password: password,
                login_type: inputType
            })
        });

        const data = await response.json();
        
        // 根据后端响应结构判断登录结果
        if (data.code === 200) {
            // 登录成功
            showMessage('登录成功！正在跳转...', 'success');
            
            // 登录成功后跳转到首页
            setTimeout(() => {
                window.location.href = '/';
            }, 1000);
        } else if (data.code === 400) {
            // 参数验证失败
            showMessage(data.msg || '参数错误，请检查输入', 'danger');
        } else if (data.code === 500) {
            // 服务器内部错误
            showMessage(data.msg || '服务器内部错误，请稍后重试', 'danger');
        } else {
            // 其他错误
            showMessage(data.msg || '登录失败，请检查账号密码', 'danger');
        }
    } catch (error) {
        console.error('登录错误:', error);
        showMessage('网络错误，请稍后重试', 'danger');
    } finally {
        // 恢复按钮状态
        loginBtn.innerHTML = originalText;
        loginBtn.disabled = false;
    }
}

// 显示消息函数
function showMessage(message, type) {
    const loginMessage = document.getElementById('loginMessage');
    loginMessage.innerHTML = `<div class="alert alert-${type} alert-dismissible fade show" role="alert">
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    </div>`;
}