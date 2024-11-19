// js/login.js

import { generateUUID } from './utils.js';

export function initializeLogin() {
    const loginForm = document.getElementById('loginForm');
    const registerBtn = document.getElementById('registerBtn');
    const cowsayContainer = document.getElementById('cowsay-container');

    // 处理注册按钮点击
    registerBtn.addEventListener('click', () => {
        window.location.href = 'register.html';
    });

    // 处理登录表单提交
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const cowIDInput = document.getElementById('id'); // 假设输入框的 id 是 'id'
        const passwordInput = document.getElementById('password');

        const cowIDValue = cowIDInput.value.trim();
        const password = passwordInput.value;

        // 输入验证
        if (!cowIDValue || !password) {
            alert('请输入 CowID 和密码');
            return;
        }

        // 将 cowID 转换为整数
        const id = parseInt(cowIDValue, 10);
        if (isNaN(id)) {
            alert('CowID 必须是一个有效的整数');
            return;
        }

        try {
            const response = await fetch('http://localhost:8081/login', { // 修改了请求路径
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ id, password }) // 修改了请求体字段
            });

            const data = await response.json();

            // 检查响应状态和返回的 code
            if (response.ok && data.code === 0) {
                const jwtToken = data.content.token; // 提取 token
                const assignedCowID = cowIDValue; // 使用输入的 cowID

                // 存储到 sessionStorage
                sessionStorage.setItem('jwtToken', jwtToken);
                sessionStorage.setItem('CowID', assignedCowID);

                // 跳转到 main.html
                window.location.href = 'main.html';
            } else {
                // 显示后端返回的错误消息
                alert(`登录失败: ${data.content || '未知错误'}`);
            }
        } catch (error) {
            console.error('登录请求失败:', error);
            alert('登录请求失败，请稍后重试');
        }
    });
}

// 初始化登录页面
document.addEventListener('DOMContentLoaded', initializeLogin);
