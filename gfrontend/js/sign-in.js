// Registration script: collect form fields and POST JSON to local API
async function registerUser() {
    // Removed event.preventDefault() as it's unnecessary for type="button"
    const username = document.getElementById('reg-username');
    const email = document.getElementById('reg-email');
    const phone = document.getElementById('reg-phone');
    const password = document.getElementById('reg-password');
    const msgEl = document.getElementById('registerMessage');

    // Basic client-side validation
    if (!username || !email || !phone || !password) {
        console.error('Form elements not found');
        return;
    }

    const payload = {
        username: username.value.trim(),
        password: password.value,
        email: email.value.trim(),
        phone: phone.value.trim()
    };

    if (!payload.username || !payload.password || !payload.email || !payload.phone) {
        msgEl.textContent = 'Please fill out all fields.';
        msgEl.className = 'text-danger';
        return;
    }

    msgEl.textContent = 'Registering...';
    msgEl.className = 'text-secondary';

    try {
        const resp = await fetch('/api/user/register', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        const result = await resp.json();

        if (result.code === 200) {
            msgEl.textContent = 'Registration successful!';
            msgEl.style.color = 'green';
            // Optionally clear form
            username.value = '';
            email.value = '';
            phone.value = '';
            password.value = '';
            // Redirect to login page
            setTimeout(() => {
                window.location.href = '/login';
            }, 2000); // Redirect after 2 seconds
        } else {
            msgEl.textContent = result.msg || 'Registration failed. Please try again.';
            msgEl.style.color = 'red';
            console.error('Register failed', resp.status, result.msg);
        }
    } catch (err) {
        msgEl.textContent = 'An error occurred. Please try again later.';
        msgEl.style.color = 'red';
        console.error('Network error', err);
    }
}

// Export for debugging in browsers that don't allow direct access to functions in modules
window.registerUser = registerUser;
