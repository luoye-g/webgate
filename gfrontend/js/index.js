// 获取Cookie的函数
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// 检查用户登录状态
async function checkUserLogin() {
    const userSession = getCookie('user_session');
    const loginButtonsCenter = document.querySelector('.login-buttons-center');
    const userInfoContainer = document.getElementById('user-info-container');

    if (!userSession) {
        // 未登录，显示中间位置的登录注册按钮
        if (loginButtonsCenter) loginButtonsCenter.style.display = 'block';
        if (userInfoContainer) userInfoContainer.style.display = 'none';
        return;
    }

    // 已登录，获取用户信息
    try {
        const response = await fetch('/api/user/info', {
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            },
            credentials: 'include' // 包含cookie
        });

        const data = await response.json();
        
        if (data.code === 200) {
            const userInfo = data.data;
            
            // 隐藏中间位置的登录注册按钮
            if (loginButtonsCenter) loginButtonsCenter.style.display = 'none';
            
            // 显示用户信息
            if (userInfoContainer) {
                const userName = userInfo.user_name || '用户';
                const initials = userName.charAt(0).toUpperCase();
                userInfoContainer.style.display = 'block';
                userInfoContainer.innerHTML = `
                    <div class="user-chip" id="userChip">
                        <div class="user-chip-trigger" onclick="toggleUserDropdown(event)">
                            <div class="chip-avatar-wrap">
                                <div class="chip-avatar">${initials}</div>
                                <div class="chip-status"></div>
                            </div>
                            <span class="chip-name">${userName}</span>
                            <span class="chip-arrow"><i class="fas fa-chevron-down"></i></span>
                        </div>
                        <div class="user-dropdown">
                            <div class="dropdown-header">
                                <div class="dropdown-header-name">${userName}</div>
                                <div class="dropdown-header-role">Blogger</div>
                            </div>
                            <div class="dropdown-menu-list">
                                <a class="dropdown-menu-item" href="/profile">
                                    <i class="fas fa-user"></i> 个人中心
                                </a>
                                <a class="dropdown-menu-item" href="/blog-create">
                                    <i class="fas fa-feather-alt"></i> 写文章
                                </a>
                                <div class="dropdown-divider"></div>
                                <button class="dropdown-menu-item" onclick="switchAccount()">
                                    <i class="fas fa-exchange-alt"></i> 切换账户
                                </button>
                                <button class="dropdown-menu-item danger" onclick="logout()">
                                    <i class="fas fa-sign-out-alt"></i> 退出登录
                                </button>
                            </div>
                        </div>
                    </div>
                `;
            }
        } else {
            // 获取用户信息失败，显示中间位置的登录按钮
            if (loginButtonsCenter) loginButtonsCenter.style.display = 'block';
            if (userInfoContainer) userInfoContainer.style.display = 'none';
        }
    } catch (error) {
        console.error('获取用户信息失败:', error);
        if (loginButtonsCenter) loginButtonsCenter.style.display = 'block';
        if (userInfoContainer) userInfoContainer.style.display = 'none';
    }
}

// Toggle user dropdown menu
function toggleUserDropdown(event) {
    event.stopPropagation();
    const chip = document.getElementById('userChip');
    if (chip) {
        chip.classList.toggle('open');
    }
}

// Close dropdown when clicking outside
document.addEventListener('click', function(event) {
    const chip = document.getElementById('userChip');
    if (chip && !chip.contains(event.target)) {
        chip.classList.remove('open');
    }
});

// Close dropdown on Escape key
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        const chip = document.getElementById('userChip');
        if (chip) chip.classList.remove('open');
    }
});

// 切换账户函数
function switchAccount() {
    window.location.href = '/login';
}

// 退出登录函数
async function logout() {
    try {
        const response = await fetch('/api/user/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include' // 包含cookie
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            // 清除cookie
            document.cookie = 'user_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            
            // 跳转到首页
            window.location.href = '/';
        } else {
            alert('退出登录失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('退出登录错误:', error);
        alert('退出登录失败，请重试');
    }
}

// 跳转到博客创建页面（需要登录检查）
function goToBlogCreate() {
    const userSession = getCookie('user_session');
    
    if (!userSession) {
        // 未登录，跳转到登录页面
        alert('请先登录后再写博客');
        window.location.href = '/login';
        return;
    }
    
    // 已登录，跳转到博客创建页面
    window.location.href = '/blog-create';
}

// 加载最新博客列表
async function loadLatestBlogs() {
    const blogLoading = document.getElementById('blogLoading');
    const blogList = document.getElementById('blogList');
    const blogEmpty = document.getElementById('blogEmpty');
    const blogError = document.getElementById('blogError');
    
    // 显示加载状态
    blogLoading.style.display = 'block';
    blogList.style.display = 'none';
    blogEmpty.style.display = 'none';
    blogError.style.display = 'none';
    
    try {
        const response = await fetch('/api/blogs?page=1&page_size=4', {
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            const blogs = data.data.blogs || [];
            
            if (blogs.length === 0) {
                // 显示空状态
                blogLoading.style.display = 'none';
                blogEmpty.style.display = 'block';
                return;
            }
            
            // 渲染博客列表 - 现代化简洁卡片
            let html = '';
            blogs.forEach(blog => {
                // 截取内容摘要（更长，充分利用空间）
                const excerpt = blog.content ? 
                    (blog.content.length > 150 ? blog.content.substring(0, 150) + '...' : blog.content) : 
                    '暂无内容';
                
                // 格式化日期
                const publishedAt = blog.published_at ? 
                    new Date(blog.published_at).toLocaleDateString('zh-CN') : 
                    (blog.created_at ? new Date(blog.created_at).toLocaleDateString('zh-CN') : '未知');
                
                html += `
                    <div class="col-12 mb-4">
                        <div class="blog-card-full">
                            <div class="card-header">
                                <h3 class="card-title">${blog.title || '无标题'}</h3>
                                <div class="meta-info">
                                    <span class="author">👤 ${blog.author_name || '未知'}</span>
                                    <span class="date">📅 ${publishedAt}</span>
                                    <span class="views">👁 ${blog.view_count || 0}</span>
                                </div>
                            </div>
                            <div class="card-body">
                                <p class="excerpt">${excerpt}</p>
                            </div>
                            <div class="card-footer">
                                <a href="/blog-public-view/${blog.id}" class="read-more">阅读全文 →</a>
                            </div>
                        </div>
                    </div>
                `;
            });
            
            blogList.innerHTML = html;
            blogLoading.style.display = 'none';
            blogList.style.display = 'block';
        } else {
            // 显示错误状态
            blogLoading.style.display = 'none';
            blogError.style.display = 'block';
            document.getElementById('blogErrorMsg').textContent = data.msg || '未知错误';
        }
    } catch (error) {
        console.error('加载博客列表失败:', error);
        // 显示错误状态
        blogLoading.style.display = 'none';
        blogError.style.display = 'block';
        document.getElementById('blogErrorMsg').textContent = error.message || '网络错误，请检查连接';
    }
}

// 页面加载完成后检查用户状态和加载博客列表
document.addEventListener('DOMContentLoaded', function() {
    checkUserLogin();
    loadLatestBlogs();
});