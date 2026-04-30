// 获取Cookie的函数
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// 当前页码和每页数量
let currentPage = 1;
let pageSize = 10;
let totalPages = 1;

// 加载博客列表
async function loadBlogList(page = 1) {
    const userSession = getCookie('user_session');

    // 检查登录状态
    if (!userSession) {
        alert('请先登录');
        window.location.href = '/login';
        return;
    }

    currentPage = page;
    
    try {
        const response = await fetch(`/api/my-blogs?page=${currentPage}&page_size=${pageSize}`, {
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            },
            credentials: 'include' // 包含cookie
        });

        const data = await response.json();
        
        if (data.code === 200) {
            const blogs = data.data.blogs || [];
            const total = data.data.total || 0;
            
            // 计算总页数
            totalPages = Math.ceil(total / pageSize);
            
            // 渲染博客列表
            renderBlogList(blogs);
            
            // 渲染分页控件
            renderPagination(totalPages);
        } else {
            alert('获取博客列表失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('获取博客列表失败:', error);
        alert('获取博客列表失败，请重试');
    }
}

// 渲染博客列表
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
    blogs.forEach(blog => {
        const statusBadge = blog.status === 1 
            ? '<span class="status-badge status-published">已发布</span>'
            : '<span class="status-badge status-draft">草稿</span>';
        
        const publishButton = blog.status === 0 
            ? `<button class="btn btn-action btn-publish" onclick="publishBlog(${blog.id})">发布</button>`
            : '';
        
        const tags = blog.tags ? blog.tags.split(',').map(tag => 
            `<span class="blog-tag">${tag.trim()}</span>`
        ).join('') : '';
        
        // 格式化日期
        const createdAt = new Date(blog.created_at).toLocaleString('zh-CN');
        const updatedAt = new Date(blog.updated_at).toLocaleString('zh-CN');
        
        html += `
            <div class="blog-card">
                <div class="d-flex justify-content-between align-items-start mb-3">
                    <h3 class="blog-card-title">${blog.title}</h3>
                    ${statusBadge}
                </div>
                
                <div class="blog-card-meta">
                    <span>创建于：${createdAt}</span>
                    <span class="mx-2">|</span>
                    <span>更新于：${updatedAt}</span>
                    <span class="mx-2">|</span>
                    <span>浏览数：${blog.view_count}</span>
                </div>
                
                ${blog.category ? `<div class="mb-2"><strong>分类：</strong>${blog.category}</div>` : ''}
                
                ${tags ? `<div class="blog-card-tags">${tags}</div>` : ''}
                
                <div class="blog-card-content">
                    ${blog.content}
                </div>
                
                <div class="blog-actions">
                    <button class="btn btn-action btn-view" onclick="viewBlog(${blog.id})">查看</button>
                    <button class="btn btn-action btn-edit" onclick="editBlog(${blog.id})">编辑</button>
                    ${publishButton}
                    <button class="btn btn-action btn-delete" onclick="deleteBlog(${blog.id})">删除</button>
                </div>
            </div>
        `;
    });
    
    blogListContainer.innerHTML = html;
}

// 渲染分页控件
function renderPagination(totalPages) {
    const paginationContainer = document.getElementById('pagination');
    
    if (totalPages <= 1) {
        paginationContainer.innerHTML = '';
        return;
    }
    
    let html = '<nav><ul class="pagination">';
    
    // 上一页按钮
    html += `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" onclick="loadBlogList(${currentPage - 1})">上一页</a>
        </li>
    `;
    
    // 页码按钮
    for (let i = 1; i <= totalPages; i++) {
        if (i === 1 || i === totalPages || (i >= currentPage - 1 && i <= currentPage + 1)) {
            html += `
                <li class="page-item ${i === currentPage ? 'active' : ''}">
                    <a class="page-link" href="#" onclick="loadBlogList(${i})">${i}</a>
                </li>
            `;
        } else if (i === currentPage - 2 || i === currentPage + 2) {
            html += '<li class="page-item disabled"><span class="page-link">...</span></li>';
        }
    }
    
    // 下一页按钮
    html += `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" onclick="loadBlogList(${currentPage + 1})">下一页</a>
        </li>
    `;
    
    html += '</ul></nav>';
    paginationContainer.innerHTML = html;
}

// 创建博客
function createBlog() {
window.open('/blog-editor', '_blank', 'noopener');
}

// 查看博客
function viewBlog(blogId) {
    window.location.href = `/blog-view/${blogId}`;
}

// 编辑博客
function editBlog(blogId) {
window.open(`/blog-editor/${blogId}`, '_blank', 'noopener');
}

// 删除博客
async function deleteBlog(blogId) {
    if (!confirm('确定要删除这篇博客吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/blogs/${blogId}`, {
            method: 'DELETE',
            headers: {
                'Accept': 'application/json'
            },
            credentials: 'include' // 包含cookie
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            alert('删除成功');
            loadBlogList(currentPage);
        } else {
            alert('删除失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('删除博客失败:', error);
        alert('删除博客失败，请重试');
    }
}

// 发布博客
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
            credentials: 'include', // 包含cookie
            body: JSON.stringify({
                blog_id: blogId
            })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            alert('发布成功');
            loadBlogList(currentPage);
        } else {
            alert('发布失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('发布博客失败:', error);
        alert('发布博客失败，请重试');
    }
}

// 页面加载完成后加载博客列表
document.addEventListener('DOMContentLoaded', function() {
    loadBlogList();
});
