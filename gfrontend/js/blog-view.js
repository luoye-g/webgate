// 获取Cookie的函数
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// 当前博客ID
let currentBlogId = null;

// 检查登录状态
function checkLogin() {
    const userSession = getCookie('user_session');
    if (!userSession) {
        alert('请先登录');
        window.location.href = '/login';
        return false;
    }
    return true;
}

// 获取博客详情
async function getBlogDetail(blogId) {
    // 这里使用 my-blogs 接口
    const response = await fetch(`/api/my-blogs?page=1&page_size=100`, {
        method: 'GET',
        headers: {
            'Accept': 'application/json'
        },
        credentials: 'include' // 包含cookie
    });
    
    const data = await response.json();
    
    if (data.code === 200) {
        const blogs = data.data.blogs || [];
        const blog = blogs.find(b => b.id === blogId);
        return blog || null;
    }
    
    return null;
}

// 加载博客信息
async function loadBlogInfo() {
    // 从URL中获取博客ID
    const pathParts = window.location.pathname.split('/');
    currentBlogId = parseInt(pathParts[pathParts.length - 1]);
    
    if (!currentBlogId || isNaN(currentBlogId)) {
        alert('无效的博客ID');
        window.location.href = '/blog-list';
        return;
    }
    
    try {
        const blog = await getBlogDetail(currentBlogId);
        
        if (!blog) {
            alert('博客不存在或已被删除');
            window.location.href = '/blog-list';
            return;
        }
        
        // 填充博客内容
        document.getElementById('blogTitle').textContent = blog.title || '-';
        
        // 设置状态徽章
        const statusBadge = document.getElementById('statusBadge');
        if (blog.status === 1) {
            statusBadge.className = 'status-badge status-published';
            statusBadge.textContent = '已发布';
        } else {
            statusBadge.className = 'status-badge status-draft';
            statusBadge.textContent = '草稿';
        }
        
        // 填充元信息
        document.getElementById('authorName').textContent = blog.author_name || '未知';
        document.getElementById('createdAt').textContent = new Date(blog.created_at).toLocaleString('zh-CN');
        document.getElementById('updatedAt').textContent = new Date(blog.updated_at).toLocaleString('zh-CN');
        document.getElementById('viewCount').textContent = blog.view_count || 0;
        
        // 分类
        if (blog.category) {
            document.getElementById('category').textContent = blog.category;
            document.getElementById('categoryContainer').style.display = 'block';
        }
        
        // 标签
        if (blog.tags) {
            const tags = blog.tags.split(',').map(tag => 
                `<span class="blog-tag">${tag.trim()}</span>`
            ).join('');
            document.getElementById('tags').innerHTML = tags;
            document.getElementById('tagsContainer').style.display = 'block';
        }
        
        // 博客内容 - 使用 Markdown 渲染
        const blogContentElement = document.getElementById('blogContent');
        if (blog.content) {
            const options = {
                breaks: true,
                gfm: true,
                sanitize: false
            };
            blogContentElement.innerHTML = marked.parse(blog.content, options);
        } else {
            blogContentElement.textContent = '-';
        }
        
        // 显示内容
        document.getElementById('loading').style.display = 'none';
        document.getElementById('blogContentContainer').style.display = 'block';
    } catch (error) {
        console.error('加载博客信息失败:', error);
        alert('加载博客信息失败，请重试');
        window.location.href = '/blog-list';
    }
}

// 编辑博客
function editBlog() {
    window.location.href = `/blog-edit/${currentBlogId}`;
}

// 删除博客
async function deleteBlog() {
    if (!confirm('确定要删除这篇博客吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/blogs/${currentBlogId}`, {
            method: 'DELETE',
            headers: {
                'Accept': 'application/json'
            },
            credentials: 'include' // 包含cookie
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            alert('删除成功');
            window.location.href = '/blog-list';
        } else {
            alert('删除失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('删除博客失败:', error);
        alert('删除博客失败，请重试');
    }
}

// 返回博客列表
function goBack() {
    window.location.href = '/blog-list';
}

// 页面加载完成后加载博客信息
document.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    if (!checkLogin()) {
        return;
    }
    
    // 加载博客信息
    loadBlogInfo();
});
