// 获取Cookie的函数
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

// Vditor 编辑器实例
let vditorInstance = null;

// 临时存储博客信息
let blogInfo = {
    category: '',
    tags: ''
};

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

// 创建博客
async function createBlog(title, content, category, tags, status, scope) {
    const response = await fetch('/api/blogs', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        credentials: 'include', // 包含cookie
        body: JSON.stringify({
            title: title,
            content: content,
            category: category || '',
            tags: tags || '',
            status: status, // 0:草稿, 1:发布
            scope: scope // 0:私密, 1:公开
        })
    });
    
    const data = await response.json();
    return data;
}

// 初始化Markdown编辑器
function initMarkdownEditor() {
    console.log('初始化Vditor编辑器...');

    return new Promise((resolve, reject) => {
        if (typeof Vditor !== 'undefined') {
            console.log('Vditor 库已加载');
            // 等待 DOM 加载完成
            setTimeout(() => {
                console.log('开始创建 Vditor 实例');
                const vditorElement = document.getElementById('vditor');
                if (!vditorElement) {
                    console.error('未找到元素 #vditor');
                    reject(new Error('未找到元素 #vditor'));
                    return;
                }

                vditorInstance = new Vditor('vditor', {
                    cdn: '/node_modules/vditor',
                    mode: 'sv', // 分屏预览模式
                    height: '100%',
                    placeholder: '开始书写你的故事…\n\n支持 Markdown 语法，尽情创作吧 ✨',
                    theme: 'dark',
                    icon: 'ant',
                    toolbar: [
                        'emoji',
                        'headings',
                        'bold',
                        'italic',
                        'strike',
                        'link',
                        '|',
                        'list',
                        'ordered-list',
                        'check',
                        'outdent',
                        'indent',
                        '|',
                        'quote',
                        'line',
                        'code',
                        'inline-code',
                        'insert-before',
                        'insert-after',
                        '|',
                        'table',
                        'upload',
                        '|',
                        'undo',
                        'redo',
                        '|',
                        'fullscreen',
                        'edit-mode',
                        'preview',
                        'help'
                    ],
                    toolbarConfig: {
                        pin: true
                    },
                    cache: {
                        enable: false
                    },
                    preview: {
                        theme: {
                            current: 'dark'
                        },
                        hljs: {
                            style: 'native'
                        }
                    },
                    input: (value) => {
                        // 输入回调，可以用于实时处理
                    },
                    after: () => {
                        console.log('Vditor 实例创建完成');
                        // Apply dark content theme
                        if (vditorInstance) {
                            vditorInstance.setTheme('dark', 'dark', 'native');
                        }
                        // Setup sync scroll between editor and preview
                        setupSyncScroll();
                        resolve();
                    }
                });
            }, 100);
        } else {
            console.error('Vditor 库未加载，请检查依赖');
            reject(new Error('Vditor 库未加载'));
        }
    });
}

// Setup sync scroll between editor and preview panels
function setupSyncScroll() {
    setTimeout(() => {
        const vditorEl = document.querySelector('.vditor');
        if (!vditorEl) return;

        // For sv (split-view) mode: sync .vditor-sv and .vditor-content .vditor-preview
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

            console.log('Sync scroll setup complete (sv mode)');
            return;
        }

        // Fallback: find any two scrollable panels in vditor-content
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

            console.log('Sync scroll setup complete (fallback mode)');
        }
    }, 500);
}

// 显示保存模态框
function saveWithOptions() {
    // 验证标题和内容
    const title = document.getElementById('blogTitle').value.trim();
    const content = vditorInstance ? vditorInstance.getValue() : '';

    if (!title) {
        alert('请输入博客标题');
        return;
    }

    if (!content) {
        alert('请输入博客内容');
        return;
    }

    // 填充模态框
    document.getElementById('modalCategory').value = blogInfo.category;
    document.getElementById('modalTags').value = blogInfo.tags;

    // 显示模态框
    document.getElementById('saveModal').classList.add('show');
}

// 关闭保存模态框
function closeSaveModal() {
    document.getElementById('saveModal').classList.remove('show');
}

// 保存为草稿
async function saveAsDraft() {
    if (!checkLogin()) {
        return;
    }
    
    const title = document.getElementById('blogTitle').value.trim();
    const content = vditorInstance ? vditorInstance.getValue() : '';
    const category = document.getElementById('modalCategory').value.trim();
    const tags = document.getElementById('modalTags').value.trim();
    const scope = parseInt(document.getElementById('modalScope').value, 10);
    
    // 验证表单
    if (!title) {
        alert('请输入博客标题');
        return;
    }
    
    if (!content) {
        alert('请输入博客内容');
        return;
    }
    
    try {
        const data = await createBlog(title, content, category, tags, 0, scope);
        
        if (data.code === 200) {
            alert('博客创建成功！');
            closeSaveModal();
            window.location.href = '/blog-list';
        } else {
            alert('创建失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('创建博客失败:', error);
        alert('创建博客失败，请重试');
    }
}

// 发布博客
async function publishBlog() {
    if (!checkLogin()) {
        return;
    }
    
    const title = document.getElementById('blogTitle').value.trim();
    const content = vditorInstance ? vditorInstance.getValue() : '';
    const category = document.getElementById('modalCategory').value.trim();
    const tags = document.getElementById('modalTags').value.trim();
    const scope = parseInt(document.getElementById('modalScope').value, 10);
    
    // 验证表单
    if (!title) {
        alert('请输入博客标题');
        return;
    }
    
    if (!content) {
        alert('请输入博客内容');
        return;
    }
    
    try {
        const data = await createBlog(title, content, category, tags, 1, scope);
        
        if (data.code === 200) {
            alert('博客发布成功！');
            closeSaveModal();
            window.location.href = '/blog-list';
        } else {
            alert('发布失败：' + (data.msg || '未知错误'));
        }
    } catch (error) {
        console.error('发布博客失败:', error);
        alert('发布博客失败，请重试');
    }
}

// 返回博客列表
function goBack() {
    if (confirm('确定要返回吗？未保存的修改将丢失。')) {
        window.location.href = '/blog-list';
    }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    if (!checkLogin()) {
        return;
    }

    // 初始化Markdown编辑器
    initMarkdownEditor().then(() => {
        console.log('编辑器初始化完成');
    }).catch(error => {
        console.error('初始化编辑器失败:', error);
        alert('初始化编辑器失败，请刷新页面重试');
    });

    // 绑定表单提交事件
    document.getElementById('blogForm').addEventListener('submit', function(event) {
        event.preventDefault();
        saveWithOptions();
    });
});
