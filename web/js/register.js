// js/register.js

export function initializeRegister() {
    const registerForm = document.getElementById('registerForm');

    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value;

        // 根据需求，收集其他必要的信息

        if (!username || !password) {
            alert('请输入用户名和密码');
            return;
        }

        // 可选：添加更多的输入合法性验证

        try {
            const response = await fetch('http://localhost:8081/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                const data = await response.json();

                alert(`注册成功！您的 CowID 是: ${data.content.id}`);
                // 跳转回登录页
                window.location.href = 'login.html';
            } else {
                const errorData = await response.json();
                alert(`注册失败: ${errorData.message || '未知错误'}`);
            }
        } catch (error) {
            console.error('注册请求失败:', error);
            alert('注册请求失败，请稍后重试');
        }
    });
}

// 初始化注册页面
document.addEventListener('DOMContentLoaded', initializeRegister);
