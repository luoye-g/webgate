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
                                <a class="dropdown-menu-item" href="/blog-editor" target="_blank" rel="noopener">
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
    
    // 已登录，新标签页打开博客创建页面
    window.open('/blog-editor', '_blank', 'noopener');
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
        const response = await fetch('/api/blogs?page=1&page_size=6', {
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
            
            // 渲染博客列表 - 九宫格卡片
            const escapeHtml = (s) => String(s == null ? '' : s)
                .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
                .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
            const stripMarkup = (s) => String(s == null ? '' : s)
                .replace(/```[\s\S]*?```/g, ' ')
                .replace(/`[^`]*`/g, ' ')
                .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
                .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
                .replace(/<[^>]+>/g, ' ')
                .replace(/[#>*_~\-]+/g, ' ')
                .replace(/\s+/g, ' ')
                .trim();

            const GRID_SIZE = 6;
            let html = '';
            blogs.slice(0, GRID_SIZE).forEach((blog, idx) => {
                const title = blog.title || '无标题';
                const author = blog.author_name || '未知';
                const plain = stripMarkup(blog.content);
                const excerpt = plain
                    ? (plain.length > 90 ? plain.substring(0, 90) + '…' : plain)
                    : '暂无内容';

                // 格式化日期
                const rawDate = blog.published_at || blog.created_at;
                const publishedAt = rawDate
                    ? new Date(rawDate).toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' })
                    : '未知';

                // 字母徽章（取标题首字）
                const firstChar = (title.trim().charAt(0) || '·').toUpperCase();
                const readMinutes = Math.max(1, Math.round((plain.length || 0) / 400));

                html += `
                    <article class="post-card" style="animation-delay:${idx * 50}ms">
                        <a class="post-card-link" href="/blog-view/${blog.id}" aria-label="${escapeHtml(title)}">
                            <div class="post-card-head">
                                <div class="post-badge" aria-hidden="true">${escapeHtml(firstChar)}</div>
                                <span class="post-views"><i class="far fa-eye"></i>${blog.view_count || 0}</span>
                            </div>
                            <h3 class="post-title">${escapeHtml(title)}</h3>
                            <p class="post-excerpt">${escapeHtml(excerpt)}</p>
                            <div class="post-foot">
                                <div class="post-meta">
                                    <span class="post-meta-item" title="${escapeHtml(author)}"><i class="far fa-user"></i>${escapeHtml(author)}</span>
                                    <span class="post-meta-dot"></span>
                                    <span class="post-meta-item"><i class="far fa-calendar"></i>${escapeHtml(publishedAt)}</span>
                                </div>
                                <span class="post-cta" aria-hidden="true"><i class="fas fa-arrow-right"></i></span>
                            </div>
                        </a>
                    </article>
                `;
            });

            // 不足 9 个用占位卡补满九宫格
            const placeholders = Math.max(0, GRID_SIZE - blogs.length);
            for (let i = 0; i < placeholders; i++) {
                html += `
                    <div class="post-card post-card--placeholder" aria-hidden="true" style="animation-delay:${(blogs.length + i) * 50}ms">
                        <div class="placeholder-inner">
                            <div class="placeholder-dot"></div>
                            <div class="placeholder-dot"></div>
                            <div class="placeholder-dot"></div>
                        </div>
                    </div>
                `;
            }
            
            blogList.innerHTML = html;
            blogLoading.style.display = 'none';
            blogList.style.display = 'grid';
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