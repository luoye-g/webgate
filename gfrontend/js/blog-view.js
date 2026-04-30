// 获取Cookie
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// 是否从"我的博客列表"页面进来（影响返回按钮文本/跳转目标）
// 通过 referrer 判断：从 /blog-list 来才视为来自列表；其它情况（首页、直接打开等）都视为访客浏览
const fromBlogList = (() => {
    try {
        if (!document.referrer) return false;
        const ref = new URL(document.referrer);
        return ref.origin === window.location.origin && ref.pathname.startsWith('/blog-list');
    } catch (e) {
        return false;
    }
})();
const isLoggedIn = !!getCookie('user_session');

// 当前博客ID
let currentBlogId = null;
let currentBlog = null;

// 获取当前登录用户信息（用于判断是否为作者本人）
async function getCurrentUser() {
    if (!isLoggedIn) return null;
    try {
        const resp = await fetch('/api/user/info', {
            method: 'GET',
            headers: { 'Accept': 'application/json' },
            credentials: 'include'
        });
        const data = await resp.json();
        if (data.code === 200) return data.data;
    } catch (e) {
        console.warn('获取用户信息失败:', e);
    }
    return null;
}

// 获取博客详情（统一走公开接口）
async function getBlogDetail(blogId) {
    const response = await fetch(`/api/blogs/${blogId}`, {
        method: 'GET',
        headers: { 'Accept': 'application/json' },
        credentials: 'include'
    });
    const data = await response.json();
    if (data.code === 200) return data.data;
    return null;
}

// 返回按钮目标：从博客列表来就回列表，否则回首页
function goBack() {
    window.location.href = fromBlogList ? '/blog-list' : '/';
}

// 生成作者头像的字母占位
function makeAvatarChar(name) {
    if (!name) return '·';
    const trimmed = String(name).trim();
    return trimmed ? trimmed.charAt(0).toUpperCase() : '·';
}

// 格式化时间
function formatDateTime(ts) {
    if (!ts) return '-';
    const d = new Date(ts);
    if (isNaN(d.getTime())) return '-';
    return d.toLocaleString('zh-CN', {
        year: 'numeric', month: '2-digit', day: '2-digit',
        hour: '2-digit', minute: '2-digit'
    });
}
function formatDateShort(ts) {
    if (!ts) return '-';
    const d = new Date(ts);
    if (isNaN(d.getTime())) return '-';
    return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' });
}

// 加载博客信息
async function loadBlogInfo() {
    const pathParts = window.location.pathname.split('/');
    currentBlogId = parseInt(pathParts[pathParts.length - 1]);

    if (!currentBlogId || isNaN(currentBlogId)) {
        alert('无效的博客ID');
        goBack();
        return;
    }

    // 配置返回按钮
    const backText = document.getElementById('backText');
    if (backText) backText.textContent = fromBlogList ? '返回列表' : '返回首页';
    const backLink = document.getElementById('backLink');
    if (backLink) {
        backLink.setAttribute('href', fromBlogList ? '/blog-list' : '/');
    }

    try {
        // 并行获取博客和用户信息
        const [blog, user] = await Promise.all([
            getBlogDetail(currentBlogId),
            getCurrentUser()
        ]);

        if (!blog) {
            alert('博客不存在或已被删除');
            goBack();
            return;
        }
        currentBlog = blog;

        // ========== 标题 ==========
        const title = blog.title || '-';
        document.getElementById('blogTitle').textContent = title;
        document.getElementById('topBarTitle').textContent = title;
        document.title = `${title} - 我的博客`;

        // ========== 标题下方副信息 ==========
        document.getElementById('kickerAuthor').textContent = blog.author_name || '未知';
        document.getElementById('kickerDate').textContent = formatDateShort(blog.created_at);
        document.getElementById('kickerViews').textContent = blog.view_count || 0;

        // ========== 状态徽章 ==========
        const statusBadge = document.getElementById('statusBadge');
        if (blog.status === 1) {
            statusBadge.className = 'status-badge status-published';
            statusBadge.textContent = '已发布';
        } else {
            statusBadge.className = 'status-badge status-draft';
            statusBadge.textContent = '草稿';
        }

        // ========== 作者卡 ==========
        document.getElementById('authorName').textContent = blog.author_name || '未知';
        document.getElementById('authorAvatar').textContent = makeAvatarChar(blog.author_name);

        // ========== 元信息 ==========
        document.getElementById('createdAt').textContent = formatDateTime(blog.created_at);
        document.getElementById('updatedAt').textContent = formatDateTime(blog.updated_at);
        document.getElementById('viewCount').textContent = blog.view_count || 0;

        if (blog.category) {
            document.getElementById('category').textContent = blog.category;
            document.getElementById('categoryRow').style.display = '';
        }

        // ========== 标签 ==========
        if (blog.tags) {
            const tagsHtml = blog.tags.split(',')
                .map(t => t.trim()).filter(Boolean)
                .map(tag => `<span class="blog-tag">${escapeHtml(tag)}</span>`)
                .join('');
            if (tagsHtml) {
                document.getElementById('tags').innerHTML = tagsHtml;
                document.getElementById('tagsContainer').style.display = '';
            }
        }

        // ========== 作者操作按钮（仅作者本人可见） ==========
        if (user && blog.author_id && user.id === blog.author_id) {
            document.getElementById('ownerActions').style.display = 'flex';
        }

        // ========== 正文 ==========
        const blogContentElement = document.getElementById('blogContent');
        if (blog.content) {
            const options = { breaks: true, gfm: true, sanitize: false };
            blogContentElement.innerHTML = marked.parse(blog.content, options);
            await renderMermaidBlocks(blogContentElement);
        } else {
            blogContentElement.textContent = '-';
        }

        // 显示内容（grid 双栏）
        document.getElementById('loading').style.display = 'none';
        document.getElementById('blogContentContainer').style.display = 'grid';

        // 初始化滚动相关（进度条 + 顶栏标题淡入）
        initScrollEffects();
    } catch (error) {
        console.error('加载博客信息失败:', error);
        alert('加载博客信息失败，请重试');
        goBack();
    }
}

// HTML 转义
function escapeHtml(str) {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

// 滚动效果：阅读进度 + 顶栏标题自动淡入
function initScrollEffects() {
    const progressBar = document.getElementById('readingProgress');
    const topBar = document.getElementById('topBar');
    const titleEl = document.getElementById('blogTitle');

    let ticking = false;
    function update() {
        const scrollTop = window.scrollY || document.documentElement.scrollTop;
        const scrollHeight = document.documentElement.scrollHeight - window.innerHeight;
        const percent = scrollHeight > 0 ? Math.min(100, (scrollTop / scrollHeight) * 100) : 0;
        if (progressBar) progressBar.style.width = percent + '%';

        // 当文章标题滚出顶部时，在顶栏里显示标题
        if (titleEl && topBar) {
            const titleBottom = titleEl.getBoundingClientRect().bottom;
            if (titleBottom < 60) {
                topBar.classList.add('is-scrolled');
            } else {
                topBar.classList.remove('is-scrolled');
            }
        }
        ticking = false;
    }
    function onScroll() {
        if (!ticking) {
            window.requestAnimationFrame(update);
            ticking = true;
        }
    }
    window.addEventListener('scroll', onScroll, { passive: true });
    update();
}

// 编辑博客
function editBlog() {
    window.open(`/blog-editor/${currentBlogId}`, '_blank', 'noopener');
}

// 删除博客
async function deleteBlog() {
    if (!confirm('确定要删除这篇博客吗？')) return;
    try {
        const response = await fetch(`/api/blogs/${currentBlogId}`, {
            method: 'DELETE',
            headers: { 'Accept': 'application/json' },
            credentials: 'include'
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

// ========== Mermaid 渲染 ==========
function decodeHtmlEntities(str) {
    const textarea = document.createElement('textarea');
    textarea.innerHTML = str;
    return textarea.value;
}
function extractMermaidSource(codeEl) {
    let raw = codeEl.innerText || codeEl.textContent || '';
    raw = decodeHtmlEntities(raw);
    return raw.trim();
}
let __mermaidInitialized = false;
async function renderMermaidBlocks(container) {
    if (typeof mermaid === 'undefined') {
        console.warn('Mermaid library is not loaded');
        return;
    }
    if (!__mermaidInitialized) {
        mermaid.initialize({
            startOnLoad: false,
            theme: 'default',
            securityLevel: 'loose',
            fontFamily: 'inherit'
        });
        __mermaidInitialized = true;
    }
    const codeBlocks = Array.from(container.querySelectorAll('pre code.language-mermaid, pre code[class*="mermaid"]'));
    if (codeBlocks.length === 0) return;

    const mermaidNodes = [];
    codeBlocks.forEach((codeEl) => {
        const pre = codeEl.parentElement;
        if (!pre) return;
        const source = extractMermaidSource(codeEl);
        if (!source) return;
        const mermaidDiv = document.createElement('div');
        mermaidDiv.className = 'mermaid';
        mermaidDiv.textContent = source;
        pre.replaceWith(mermaidDiv);
        mermaidNodes.push(mermaidDiv);
    });
    if (mermaidNodes.length === 0) return;

    for (let i = 0; i < mermaidNodes.length; i++) {
        const node = mermaidNodes[i];
        const source = node.textContent;
        const id = `mermaid-svg-${Date.now()}-${i}`;
        try {
            const { svg, bindFunctions } = await mermaid.render(id, source);
            node.innerHTML = svg;
            if (typeof bindFunctions === 'function') {
                bindFunctions(node);
            }
        } catch (err) {
            console.error('Mermaid render error:', err);
            const errorDiv = document.createElement('div');
            errorDiv.className = 'mermaid-error';
            errorDiv.textContent = `Mermaid 渲染失败: ${err && err.message ? err.message : err}\n\n原始内容:\n${source}`;
            node.replaceWith(errorDiv);
        }
    }
}

// 启动
document.addEventListener('DOMContentLoaded', function () {
    loadBlogInfo();
});
