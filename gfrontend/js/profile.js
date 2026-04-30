// ====== Utility Functions ======

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// ====== State ======
let currentPage = 1;
let pageSize = 8;
let totalPages = 1;
let currentFilter = 'all';
let allBlogsCache = [];

// ====== User Profile ======

async function loadUserProfile() {
    const userSession = getCookie('user_session');
    if (!userSession) {
        alert('请先登录');
        window.location.href = '/login';
        return;
    }

    try {
        const response = await fetch('/api/user/info', {
            method: 'GET',
            headers: { 'Accept': 'application/json' },
            credentials: 'include'
        });

        const data = await response.json();

        if (data.code === 200) {
            const userInfo = data.data;
            document.getElementById('profileUserName').textContent = userInfo.user_name || '用户';
            document.getElementById('profileUserNameValue').textContent = userInfo.user_name || '未设置';
            document.getElementById('profileEmail').textContent = userInfo.email || '未设置';
            document.getElementById('profilePhone').textContent = userInfo.phone || '未设置';
        } else {
            alert('获取用户信息失败：' + (data.msg || '未知错误'));
            window.location.href = '/';
        }
    } catch (error) {
        console.error('获取用户信息失败:', error);
        alert('获取用户信息失败，请重试');
        window.location.href = '/';
    }
}

// ====== Edit Nickname ======

function showNicknameEdit() {
    const currentName = document.getElementById('profileUserName').textContent;
    const input = document.getElementById('nicknameInput');
    const form = document.getElementById('nicknameEditForm');

    input.value = currentName === '用户' ? '' : currentName;
    form.classList.add('active');
    input.focus();
}

function cancelNicknameEdit() {
    const form = document.getElementById('nicknameEditForm');
    form.classList.remove('active');
}

async function saveNickname() {
    const input = document.getElementById('nicknameInput');
    const saveBtn = document.getElementById('btnNicknameSave');
    const newName = input.value.trim();

    if (!newName) {
        alert('昵称不能为空');
        input.focus();
        return;
    }

    if (newName.length > 20) {
        alert('昵称长度不能超过20个字符');
        input.focus();
        return;
    }

    saveBtn.disabled = true;
    saveBtn.textContent = '保存中...';

    try {
        const response = await fetch('/api/user/info', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ user_name: newName })
        });

        const data = await response.json();

        if (data.code === 200) {
            document.getElementById('profileUserName').textContent = newName;
            document.getElementById('profileUserNameValue').textContent = newName;
            cancelNicknameEdit();
        } else {
            alert('修改失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('修改昵称失败:', error);
        alert('修改昵称失败，请重试');
    } finally {
        saveBtn.disabled = false;
        saveBtn.textContent = '保存';
    }
}

// ====== Blog List ======

async function loadBlogList(page = 1) {
    const userSession = getCookie('user_session');
    if (!userSession) {
        alert('请先登录');
        window.location.href = '/login';
        return;
    }

    currentPage = page;

    document.getElementById('blogLoading').style.display = 'block';
    document.getElementById('blogList').innerHTML = '';
    document.getElementById('emptyState').style.display = 'none';
    document.getElementById('pagination').innerHTML = '';

    try {
        const response = await fetch(`/api/my-blogs?page=${currentPage}&page_size=${pageSize}`, {
            method: 'GET',
            headers: { 'Accept': 'application/json' },
            credentials: 'include'
        });

        const data = await response.json();

        document.getElementById('blogLoading').style.display = 'none';

        if (data.code === 200) {
            const blogs = data.data.blogs || [];
            const total = data.data.total || 0;

            totalPages = Math.ceil(total / pageSize);
            allBlogsCache = blogs;

            updateStats(blogs, total);

            const filteredBlogs = applyFilter(blogs);
            renderBlogList(filteredBlogs);
            renderPagination(totalPages);
        } else {
            alert('获取博客列表失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        document.getElementById('blogLoading').style.display = 'none';
        console.error('获取博客列表失败:', error);
        alert('获取博客列表失败，请重试');
    }
}

function updateStats(blogs, total) {
    const published = blogs.filter(b => b.status === 1).length;
    const draft = blogs.filter(b => b.status === 0).length;

    document.getElementById('statTotal').textContent = total;
    document.getElementById('statPublished').textContent = published;
    document.getElementById('statDraft').textContent = draft;

    document.getElementById('countAll').textContent = blogs.length;
    document.getElementById('countPublished').textContent = published;
    document.getElementById('countDraft').textContent = draft;
}

function applyFilter(blogs) {
    if (currentFilter === 'published') {
        return blogs.filter(b => b.status === 1);
    } else if (currentFilter === 'draft') {
        return blogs.filter(b => b.status === 0);
    }
    return blogs;
}

function filterBlogs(filter) {
    currentFilter = filter;

    document.querySelectorAll('.filter-tab').forEach(tab => {
        tab.classList.toggle('active', tab.dataset.filter === filter);
    });

    const filteredBlogs = applyFilter(allBlogsCache);
    renderBlogList(filteredBlogs);
}

function renderBlogList(blogs) {
    const blogListContainer = document.getElementById('blogList');
    const emptyState = document.getElementById('emptyState');

    if (blogs.length === 0) {
        blogListContainer.innerHTML = '';
        emptyState.style.display = 'block';
        return;
    }

    emptyState.style.display = 'none';

    let html = '';
    blogs.forEach((blog, index) => {
        const statusBadge = blog.status === 1
            ? '<span class="status-badge status-published"><i class="fas fa-check-circle"></i> 已发布</span>'
            : '<span class="status-badge status-draft"><i class="fas fa-pencil-alt"></i> 草稿</span>';

        const publishButton = blog.status === 0
            ? `<button class="btn btn-action btn-publish" onclick="publishBlog(${blog.id})"><i class="fas fa-paper-plane"></i> 发布</button>`
            : '';

        const tags = blog.tags ? blog.tags.split(',').map(tag =>
            `<span class="blog-tag">${tag.trim()}</span>`
        ).join('') : '';

        const createdAt = new Date(blog.created_at).toLocaleString('zh-CN');
        const updatedAt = new Date(blog.updated_at).toLocaleString('zh-CN');

        const contentPreview = blog.content
            ? blog.content.replace(/<[^>]*>/g, '').substring(0, 150)
            : '';

        html += `
            <div class="blog-card" style="animation: fadeSlideUp 0.5s ease ${0.05 * index}s both;">
                <div class="blog-card-header">
                    <h3 class="blog-card-title" onclick="viewBlog(${blog.id})">${blog.title}</h3>
                    ${statusBadge}
                </div>

                <div class="blog-card-meta">
                    <span><i class="far fa-clock"></i> ${createdAt}</span>
                    <span><i class="far fa-edit"></i> ${updatedAt}</span>
                    ${blog.category ? `<span><i class="fas fa-folder"></i> ${blog.category}</span>` : ''}
                </div>

                ${tags ? `<div class="blog-card-tags">${tags}</div>` : ''}

                ${contentPreview ? `<div class="blog-card-content">${contentPreview}</div>` : ''}

                <div class="blog-card-footer">
                    <div class="blog-card-stats">
                        <span class="blog-stat"><i class="far fa-eye"></i> ${blog.view_count || 0}</span>
                    </div>
                    <div class="blog-actions">
                        <button class="btn btn-action btn-view" onclick="viewBlog(${blog.id})"><i class="far fa-eye"></i> 查看</button>
                        <button class="btn btn-action btn-edit" onclick="editBlog(${blog.id})"><i class="far fa-edit"></i> 编辑</button>
                        ${publishButton}
                        <button class="btn btn-action btn-delete" onclick="deleteBlog(${blog.id})"><i class="far fa-trash-alt"></i> 删除</button>
                    </div>
                </div>
            </div>
        `;
    });

    blogListContainer.innerHTML = html;
}

function renderPagination(totalPages) {
    const paginationContainer = document.getElementById('pagination');

    if (totalPages <= 1) {
        paginationContainer.innerHTML = '';
        return;
    }

    let html = '<nav><ul class="pagination" style="list-style:none; display:flex; gap:6px; padding:0; margin:0;">';

    html += `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="javascript:void(0)" onclick="${currentPage > 1 ? `loadBlogList(${currentPage - 1})` : ''}">
                <i class="fas fa-chevron-left" style="font-size:0.75rem;"></i>
            </a>
        </li>
    `;

    for (let i = 1; i <= totalPages; i++) {
        if (i === 1 || i === totalPages || (i >= currentPage - 1 && i <= currentPage + 1)) {
            html += `
                <li class="page-item ${i === currentPage ? 'active' : ''}">
                    <a class="page-link" href="javascript:void(0)" onclick="loadBlogList(${i})">${i}</a>
                </li>
            `;
        } else if (i === currentPage - 2 || i === currentPage + 2) {
            html += '<li class="page-item disabled"><span class="page-link">...</span></li>';
        }
    }

    html += `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="javascript:void(0)" onclick="${currentPage < totalPages ? `loadBlogList(${currentPage + 1})` : ''}">
                <i class="fas fa-chevron-right" style="font-size:0.75rem;"></i>
            </a>
        </li>
    `;

    html += '</ul></nav>';
    paginationContainer.innerHTML = html;
}

// ====== Blog Actions ======

function createBlog() {
window.open('/blog-editor', '_blank', 'noopener');
}

function viewBlog(blogId) {
    window.location.href = `/blog-view/${blogId}`;
}

function editBlog(blogId) {
window.open(`/blog-editor/${blogId}`, '_blank', 'noopener');
}

async function deleteBlog(blogId) {
    if (!confirm('确定要删除这篇博客吗？此操作不可撤销。')) {
        return;
    }

    try {
        const response = await fetch(`/api/blogs/${blogId}`, {
            method: 'DELETE',
            headers: { 'Accept': 'application/json' },
            credentials: 'include'
        });

        const data = await response.json();

        if (data.code === 200) {
            loadBlogList(currentPage);
        } else {
            alert('删除失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('删除博客失败:', error);
        alert('删除博客失败，请重试');
    }
}

async function publishBlog(blogId) {
    if (!confirm('确定要发布这篇博客吗？')) {
        return;
    }

    try {
        const response = await fetch('/api/blogs/publish', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ blog_id: blogId })
        });

        const data = await response.json();

        if (data.code === 200) {
            loadBlogList(currentPage);
        } else {
            alert('发布失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('发布博客失败:', error);
        alert('发布博客失败，请重试');
    }
}

// ====== Account Actions ======

function switchAccount() {
    window.location.href = '/login';
}

async function logout() {
    try {
        const response = await fetch('/api/user/logout', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include'
        });

        const data = await response.json();

        if (data.code === 200) {
            document.cookie = 'user_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            window.location.href = '/';
        } else {
            alert('退出登录失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('退出登录错误:', error);
        alert('退出登录失败，请重试');
    }
}

// ====== Mobile Sidebar ======

function toggleSidebar() {
    const sidebar = document.getElementById('profileSidebar');
    const overlay = document.getElementById('sidebarOverlay');
    const isOpen = sidebar.classList.contains('open');

    if (isOpen) {
        sidebar.classList.remove('open');
        overlay.classList.remove('active');
    } else {
        sidebar.classList.add('open');
        overlay.classList.add('active');
    }
}

// ====== Init ======

document.addEventListener('DOMContentLoaded', function() {
    loadUserProfile();
    loadBlogList();
});
