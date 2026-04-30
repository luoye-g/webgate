// ============================================================
// blog-editor.js - 统一的博客创建/编辑页逻辑
// 根据 URL 决定模式：
//   /blog-editor            -> create 模式
//   /blog-editor/new        -> create 模式
//   /blog-editor/:id (数字) -> edit 模式
// ============================================================

// 获取Cookie
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// ============ 模式识别 ============
function detectMode() {
    const pathParts = window.location.pathname.split('/').filter(Boolean);
    // pathParts: ['blog-editor'] 或 ['blog-editor', 'new'] 或 ['blog-editor', '123']
    if (pathParts.length < 2 || pathParts[1] === 'new') {
        return { mode: 'create', blogId: null };
    }
    const id = parseInt(pathParts[1], 10);
    if (isNaN(id)) {
        return { mode: 'create', blogId: null };
    }
    return { mode: 'edit', blogId: id };
}

const { mode: pageMode, blogId: currentBlogId } = detectMode();
const isCreateMode = pageMode === 'create';

// ============ 全局状态 ============
let vditorInstance = null;
let blogInfo = {
    category: '',
    tags: '',
    scope: 1 // 默认公开
};

// ============ 登录检查 ============
function checkLogin() {
    const userSession = getCookie('user_session');
    if (!userSession) {
        alert('请先登录');
        window.location.href = '/login';
        return false;
    }
    return true;
}

// ============ UI 配置（根据模式切换文案） ============
function configureUIByMode() {
    if (isCreateMode) {
        document.title = '创建博客 - 我的博客';
        document.getElementById('saveBtnIcon').className = 'fas fa-feather-alt';
        document.getElementById('saveBtnText').textContent = '发布文章';
        document.getElementById('modalTitle').textContent = '发布设置';
        document.getElementById('draftBtnText').textContent = '存为草稿';
        document.getElementById('publishBtnText').textContent = '立即发布';
    } else {
        document.title = '编辑博客 - 我的博客';
        document.getElementById('saveBtnIcon').className = 'fas fa-save';
        document.getElementById('saveBtnText').textContent = '保存文章';
        document.getElementById('modalTitle').textContent = '保存设置';
        document.getElementById('draftBtnText').textContent = '保存草稿';
        document.getElementById('publishBtnText').textContent = '发布';
    }
}

// ============ 接口封装 ============
async function createBlog(title, content, category, tags, status, scope) {
    const response = await fetch('/api/blogs', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
            title, content,
            category: category || '',
            tags: tags || '',
            status, // 0:草稿, 1:发布
            scope   // 0:私密, 1:公开
        })
    });
    return await response.json();
}

async function updateBlog(blogId, title, content, category, tags, status, scope) {
    const response = await fetch(`/api/blogs/${blogId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
            title, content,
            category: category || '',
            tags: tags || '',
            status, scope
        })
    });
    return await response.json();
}

// edit 模式下获取博客详情（沿用原 blog-edit.js 的做法，从 my-blogs 列表中查找）
async function fetchBlogForEdit(blogId) {
    const response = await fetch(`/api/my-blogs?page=1&page_size=100`, {
        method: 'GET',
        headers: { 'Accept': 'application/json' },
        credentials: 'include'
    });
    const data = await response.json();
    if (data.code === 200) {
        const blogs = data.data.blogs || [];
        return blogs.find(b => b.id === blogId) || null;
    }
    return null;
}

// ============ Vditor 编辑器初始化 ============
function initMarkdownEditor() {
    return new Promise((resolve, reject) => {
        if (typeof Vditor === 'undefined') {
            reject(new Error('Vditor 库未加载'));
            return;
        }

        setTimeout(() => {
            const vditorElement = document.getElementById('vditor');
            if (!vditorElement) {
                reject(new Error('未找到元素 #vditor'));
                return;
            }

            const placeholder = isCreateMode
                ? '开始书写你的故事…\n\n支持 Markdown 语法，尽情创作吧 ✨'
                : '继续书写你的故事…\n\n支持 Markdown 语法，尽情创作吧 ✨';

            vditorInstance = new Vditor('vditor', {
                cdn: '/node_modules/vditor',
                mode: 'sv',
                height: '100%',
                placeholder: placeholder,
                theme: 'classic',
                icon: 'ant',
                toolbar: [
                    'emoji', 'headings', 'bold', 'italic', 'strike', 'link',
                    '|',
                    'list', 'ordered-list', 'check', 'outdent', 'indent',
                    '|',
                    'quote', 'line', 'code', 'inline-code', 'insert-before', 'insert-after',
                    '|',
                    'table', 'upload',
                    '|',
                    'undo', 'redo',
                    '|',
                    'fullscreen', 'edit-mode', 'preview', 'help'
                ],
                toolbarConfig: { pin: true },
                cache: { enable: false },
                preview: {
                    theme: { current: 'light' },
                    hljs: { style: 'github' }
                },
                input: () => {},
                after: () => {
                    if (vditorInstance) {
                        vditorInstance.setTheme('classic', 'light', 'github');
                    }
                    setupSyncScroll();
                    resolve();
                }
            });
        }, 100);
    });
}

// 编辑器 / 预览面板同步滚动
function setupSyncScroll() {
    setTimeout(() => {
        const vditorEl = document.querySelector('.vditor');
        if (!vditorEl) return;

        const editorPanel = vditorEl.querySelector('.vditor-sv .vditor-reset');
        const previewPanel = vditorEl.querySelector('.vditor-content .vditor-preview .vditor-reset');

        if (editorPanel && previewPanel) {
            let isSyncingEditor = false;
            let isSyncingPreview = false;

            editorPanel.addEventListener('scroll', () => {
                if (isSyncingEditor) { isSyncingEditor = false; return; }
                isSyncingPreview = true;
                const ratio = editorPanel.scrollTop / (editorPanel.scrollHeight - editorPanel.clientHeight || 1);
                previewPanel.scrollTop = ratio * (previewPanel.scrollHeight - previewPanel.clientHeight);
            });

            previewPanel.addEventListener('scroll', () => {
                if (isSyncingPreview) { isSyncingPreview = false; return; }
                isSyncingEditor = true;
                const ratio = previewPanel.scrollTop / (previewPanel.scrollHeight - previewPanel.clientHeight || 1);
                editorPanel.scrollTop = ratio * (editorPanel.scrollHeight - editorPanel.clientHeight);
            });
            return;
        }

        // Fallback
        const panels = vditorEl.querySelectorAll('.vditor-content > div');
        if (panels.length >= 2) {
            const leftPanel = panels[0].querySelector('.vditor-reset') || panels[0];
            const rightPanel = panels[1].querySelector('.vditor-reset') || panels[1];

            let isSyncingLeft = false;
            let isSyncingRight = false;

            leftPanel.addEventListener('scroll', () => {
                if (isSyncingLeft) { isSyncingLeft = false; return; }
                isSyncingRight = true;
                const ratio = leftPanel.scrollTop / (leftPanel.scrollHeight - leftPanel.clientHeight || 1);
                rightPanel.scrollTop = ratio * (rightPanel.scrollHeight - rightPanel.clientHeight);
            });

            rightPanel.addEventListener('scroll', () => {
                if (isSyncingRight) { isSyncingRight = false; return; }
                isSyncingLeft = true;
                const ratio = rightPanel.scrollTop / (rightPanel.scrollHeight - rightPanel.clientHeight || 1);
                leftPanel.scrollTop = ratio * (leftPanel.scrollHeight - leftPanel.clientHeight);
            });
        }
    }, 500);
}

// ============ edit 模式加载博客内容 ============
async function loadBlogInfoForEdit() {
    if (!currentBlogId || isNaN(currentBlogId)) {
        alert('无效的博客ID');
        window.location.href = '/blog-list';
        return;
    }

    try {
        const blog = await fetchBlogForEdit(currentBlogId);
        if (!blog) {
            alert('博客不存在或已被删除');
            window.location.href = '/blog-list';
            return;
        }

        // 填充表单
        document.getElementById('blogTitle').value = blog.title || '';

        blogInfo.category = blog.category || '';
        blogInfo.tags = blog.tags || '';
        blogInfo.scope = blog.scope !== undefined ? blog.scope : 1;

        if (vditorInstance) {
            vditorInstance.setValue(blog.content || '');
        }

        // 切换显示状态
        document.getElementById('loading').style.display = 'none';
        document.getElementById('titleBar').style.display = 'flex';
        document.getElementById('blogFormContainer').style.display = 'flex';
    } catch (error) {
        console.error('加载博客信息失败:', error);
        alert('加载博客信息失败，请重试');
        window.location.href = '/blog-list';
    }
}

// ============ 保存模态框 ============
function saveWithOptions() {
    const title = document.getElementById('blogTitle').value.trim();
    const content = vditorInstance ? vditorInstance.getValue() : '';

    if (!title) { alert('请输入博客标题'); return; }
    if (!content) { alert('请输入博客内容'); return; }

    // 填充模态框
    document.getElementById('modalCategory').value = blogInfo.category;
    document.getElementById('modalTags').value = blogInfo.tags;
    // 只有编辑模式下才把已有 scope 回填；创建模式保留默认值
    if (!isCreateMode) {
        document.getElementById('modalScope').value = blogInfo.scope;
    }

    document.getElementById('saveModal').classList.add('show');
}

function closeSaveModal() {
    document.getElementById('saveModal').classList.remove('show');
}

// ============ 统一的保存入口 ============
// status: 0=草稿, 1=发布
async function submitBlog(status) {
    if (!checkLogin()) return;

    const title = document.getElementById('blogTitle').value.trim();
    const content = vditorInstance ? vditorInstance.getValue() : '';
    const category = document.getElementById('modalCategory').value.trim();
    const tags = document.getElementById('modalTags').value.trim();
    const scope = parseInt(document.getElementById('modalScope').value, 10);

    if (!title) { alert('请输入博客标题'); return; }
    if (!content) { alert('请输入博客内容'); return; }

    try {
        const data = isCreateMode
            ? await createBlog(title, content, category, tags, status, scope)
            : await updateBlog(currentBlogId, title, content, category, tags, status, scope);

        if (data.code === 200) {
            const successMsg = isCreateMode
                ? (status === 1 ? '博客发布成功！' : '博客创建成功！')
                : (status === 1 ? '博客发布成功！' : '博客保存成功！');
            alert(successMsg);
            closeSaveModal();
            window.location.href = '/blog-list';
        } else {
            const failMsg = isCreateMode
                ? (status === 1 ? '发布失败' : '创建失败')
                : (status === 1 ? '发布失败' : '保存失败');
            alert(failMsg + '：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('保存博客失败:', error);
        alert('保存博客失败，请重试');
    }
}

// 供 HTML onclick 调用
function saveAsDraft() { return submitBlog(0); }
function publishBlog() { return submitBlog(1); }

// 返回
function goBack() {
    if (confirm('确定要返回吗？未保存的修改将丢失。')) {
        window.location.href = '/blog-list';
    }
}

// ============ 启动 ============
document.addEventListener('DOMContentLoaded', function () {
    if (!checkLogin()) return;

    // 配置 UI 文案
    configureUIByMode();

    // create 模式：直接显示编辑器；edit 模式：先显示 Loading
    if (isCreateMode) {
        document.getElementById('loading').style.display = 'none';
        document.getElementById('titleBar').style.display = 'flex';
        document.getElementById('blogFormContainer').style.display = 'flex';
    } else {
        document.getElementById('loading').style.display = 'flex';
        document.getElementById('titleBar').style.display = 'none';
        document.getElementById('blogFormContainer').style.display = 'none';
    }

    // 初始化编辑器
    initMarkdownEditor()
        .then(() => {
            if (!isCreateMode) {
                return loadBlogInfoForEdit();
            }
        })
        .catch(error => {
            console.error('初始化编辑器失败:', error);
            alert('初始化编辑器失败，请刷新页面重试');
        });

    // 表单回车提交
    document.getElementById('blogForm').addEventListener('submit', function (event) {
        event.preventDefault();
        saveWithOptions();
    });
});
