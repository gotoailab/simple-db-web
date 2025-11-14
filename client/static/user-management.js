(function() {
    // 等待页面加载完成
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initUserManagement);
    } else {
        initUserManagement();
    }

    function initUserManagement() {
        // 获取当前用户信息
        fetch('/api/auth/current')
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    addUserManagementUI(data.data);
                }
            })
            .catch(err => console.error('Failed to get user info:', err));
    }

    function addUserManagementUI(currentUser) {
        // 在header中添加用户菜单
        const headerContent = document.querySelector('.header-content');
        if (!headerContent) return;

        // 获取翻译函数（如果存在）
        const t = typeof window.t === 'function' ? window.t : function(key) { return key; };
        
        // 保存 currentUser 到闭包中，供后续函数使用
        const currentUserRef = currentUser;

        // 创建用户菜单容器 - 使用flex-shrink: 0 防止被压缩
        const userMenu = document.createElement('div');
        userMenu.style.cssText = 'margin-left: auto; display: flex; align-items: center; gap: 1rem; position: relative; flex-shrink: 0;';

        // 用户名显示
        const usernameSpan = document.createElement('span');
        usernameSpan.textContent = currentUser.username;
        usernameSpan.style.cssText = 'color: var(--text-primary); font-size: 0.875rem; white-space: nowrap;';
        if (currentUser.is_admin) {
            usernameSpan.textContent += ' (' + t('user.admin') + ')';
            usernameSpan.style.color = 'var(--primary-color)';
        }

        // 用户菜单按钮
        const menuButton = document.createElement('button');
        menuButton.textContent = '⚙️';
        menuButton.style.cssText = 'background: var(--surface); border: 1px solid var(--border-color); border-radius: 4px; padding: 0.5rem; cursor: pointer; color: var(--text-primary); flex-shrink: 0;';
        menuButton.title = t('user.menu');

        // 下拉菜单
        const dropdown = document.createElement('div');
        dropdown.style.cssText = 'display: none; position: absolute; top: 100%; right: 0; background: var(--surface); border: 1px solid var(--border-color); border-radius: 4px; box-shadow: var(--shadow); margin-top: 0.5rem; min-width: 200px; z-index: 1000;';
        
        const menuItems = [
            { text: t('user.changePassword'), key: 'user.changePassword', action: showChangePasswordModal },
            { text: t('user.logout'), key: 'user.logout', action: handleLogout, style: 'color: var(--danger-color);' }
        ];

        if (currentUser.is_admin) {
            menuItems.splice(1, 0, { text: t('user.management'), key: 'user.management', action: showUserManagementModal });
        }

        menuItems.forEach(item => {
            const menuItem = document.createElement('div');
            menuItem.textContent = item.text;
            menuItem.setAttribute('data-i18n-key', item.key);
            menuItem.style.cssText = 'padding: 0.75rem 1rem; cursor: pointer; color: var(--text-primary); border-bottom: 1px solid var(--border-color); transition: background 0.2s;' + (item.style || '');
            menuItem.addEventListener('click', () => {
                item.action();
                dropdown.style.display = 'none';
            });
            menuItem.addEventListener('mouseenter', () => {
                menuItem.style.background = 'var(--surface-light)';
            });
            menuItem.addEventListener('mouseleave', () => {
                menuItem.style.background = 'transparent';
            });
            dropdown.appendChild(menuItem);
        });
        
        // 监听语言变化事件，更新菜单文本和用户名显示
        const languageChangeHandler = () => {
            menuItems.forEach((item, index) => {
                const menuItem = dropdown.children[index];
                if (menuItem && menuItem.getAttribute('data-i18n-key') === item.key) {
                    menuItem.textContent = t(item.key);
                }
            });
            // 更新用户名显示中的"Admin"文本
            if (currentUserRef.is_admin) {
                usernameSpan.textContent = currentUserRef.username + ' (' + t('user.admin') + ')';
            } else {
                usernameSpan.textContent = currentUserRef.username;
            }
        };
        window.addEventListener('languageChanged', languageChangeHandler);

        menuButton.addEventListener('click', (e) => {
            e.stopPropagation();
            dropdown.style.display = dropdown.style.display === 'none' ? 'block' : 'none';
        });

        document.addEventListener('click', (e) => {
            if (!userMenu.contains(e.target)) {
                dropdown.style.display = 'none';
            }
        });

        userMenu.appendChild(usernameSpan);
        userMenu.appendChild(menuButton);
        userMenu.appendChild(dropdown);
        headerContent.appendChild(userMenu);

        // 修改密码模态框
        function showChangePasswordModal() {
            const modal = createModal(t('user.changePassword'), 
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: block; margin-bottom: 0.5rem; color: var(--text-primary);">' + t('user.oldPassword') + '</label>' +
                    '<input type="password" id="oldPassword" style="width: 100%; padding: 0.5rem; border: 1px solid var(--border-color); border-radius: 4px; background: var(--surface); color: var(--text-primary);">' +
                '</div>' +
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: block; margin-bottom: 0.5rem; color: var(--text-primary);">' + t('user.newPassword') + '</label>' +
                    '<input type="password" id="newPassword" style="width: 100%; padding: 0.5rem; border: 1px solid var(--border-color); border-radius: 4px; background: var(--surface); color: var(--text-primary);">' +
                '</div>' +
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: block; margin-bottom: 0.5rem; color: var(--text-primary);">' + t('user.confirmPassword') + '</label>' +
                    '<input type="password" id="confirmPassword" style="width: 100%; padding: 0.5rem; border: 1px solid var(--border-color); border-radius: 4px; background: var(--surface); color: var(--text-primary);">' +
                '</div>' +
                '<div style="display: flex; gap: 0.5rem; justify-content: flex-end;">' +
                    '<button id="cancelPasswordBtn" class="btn btn-secondary">' + t('common.cancel') + '</button>' +
                    '<button id="savePasswordBtn" class="btn btn-primary">' + t('common.save') + '</button>' +
                '</div>'
            );

            document.getElementById('savePasswordBtn').addEventListener('click', async () => {
                const oldPassword = document.getElementById('oldPassword').value;
                const newPassword = document.getElementById('newPassword').value;
                const confirmPassword = document.getElementById('confirmPassword').value;

                if (!oldPassword || !newPassword || !confirmPassword) {
                    showNotification(t('user.fillAllFields'), 'error');
                    return;
                }

                if (newPassword !== confirmPassword) {
                    showNotification(t('user.passwordMismatch'), 'error');
                    return;
                }

                try {
                    const response = await fetch('/api/auth/password', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
                    });

                    const data = await response.json();
                    if (data.success) {
                        showNotification(t('user.passwordUpdated'), 'success');
                        closeModal(modal);
                    } else {
                        showNotification(data.message || t('user.passwordUpdateFailed'), 'error');
                    }
                } catch (error) {
                    showNotification('Network error: ' + error.message, 'error');
                }
            });

            document.getElementById('cancelPasswordBtn').addEventListener('click', () => {
                closeModal(modal);
            });
        }

        // 用户管理模态框（仅管理员）
        function showUserManagementModal() {
            const modal = createModal(t('user.management'), 
                '<div style="margin-bottom: 1rem;">' +
                    '<button id="addUserBtn" class="btn btn-primary" style="margin-bottom: 1rem;">+ ' + t('user.addUser') + '</button>' +
                    '<div id="usersList" style="max-height: 400px; overflow-y: auto;"></div>' +
                '</div>'
            );

            loadUsersList();

            document.getElementById('addUserBtn').addEventListener('click', () => {
                showAddUserModal();
            });
        }

        function loadUsersList() {
            fetch('/api/users')
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        displayUsersList(data.data);
                    }
                })
                .catch(err => {
                    console.error('Failed to load users:', err);
                    showNotification(t('user.loadFailed'), 'error');
                });
        }

        function displayUsersList(users) {
            const usersList = document.getElementById('usersList');
            usersList.innerHTML = '';

            users.forEach(user => {
                const userItem = document.createElement('div');
                userItem.style.cssText = 'padding: 1rem; border: 1px solid var(--border-color); border-radius: 4px; margin-bottom: 0.5rem; background: var(--surface);';
                userItem.innerHTML = 
                    '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                        '<div>' +
                            '<div style="font-weight: 600; color: var(--text-primary);">' + escapeHtml(user.username) + '</div>' +
                            '<div style="font-size: 0.875rem; color: var(--text-secondary);">' +
                                (user.is_admin ? '<span style="color: var(--primary-color);">' + t('user.admin') + '</span>' : t('user.user')) +
                            '</div>' +
                        '</div>' +
                        '<div style="display: flex; gap: 0.5rem;">' +
                            '<button class="btn btn-secondary edit-user-btn" data-id="' + user.id + '" data-username="' + escapeHtml(user.username) + '" data-is-admin="' + user.is_admin + '">' + t('common.edit') + '</button>' +
                            (user.id !== currentUserRef.id ? '<button class="btn btn-danger delete-user-btn" data-id="' + user.id + '">' + t('common.delete') + '</button>' : '') +
                        '</div>' +
                    '</div>';
                usersList.appendChild(userItem);
            });

            // 绑定编辑和删除事件
            usersList.querySelectorAll('.edit-user-btn').forEach(btn => {
                btn.addEventListener('click', () => {
                    const id = btn.dataset.id;
                    const username = btn.dataset.username;
                    const isAdmin = btn.dataset.isAdmin === 'true';
                    showEditUserModal(id, username, isAdmin);
                });
            });

            usersList.querySelectorAll('.delete-user-btn').forEach(btn => {
                btn.addEventListener('click', () => {
                    const id = btn.dataset.id;
                    if (confirm(t('user.deleteConfirm'))) {
                        deleteUser(id);
                    }
                });
            });
        }

        function showAddUserModal() {
            const modal = createModal(t('user.addUser'), 
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: block; margin-bottom: 0.5rem; color: var(--text-primary);">' + t('connection.user') + '</label>' +
                    '<input type="text" id="newUsername" style="width: 100%; padding: 0.5rem; border: 1px solid var(--border-color); border-radius: 4px; background: var(--surface); color: var(--text-primary);">' +
                '</div>' +
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: block; margin-bottom: 0.5rem; color: var(--text-primary);">' + t('connection.password') + '</label>' +
                    '<input type="password" id="newUserPassword" style="width: 100%; padding: 0.5rem; border: 1px solid var(--border-color); border-radius: 4px; background: var(--surface); color: var(--text-primary);">' +
                '</div>' +
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: flex; align-items: center; gap: 0.5rem; color: var(--text-primary);">' +
                        '<input type="checkbox" id="newUserIsAdmin">' +
                        '<span>' + t('user.admin') + '</span>' +
                    '</label>' +
                '</div>' +
                '<div style="display: flex; gap: 0.5rem; justify-content: flex-end;">' +
                    '<button id="cancelAddUserBtn" class="btn btn-secondary">' + t('common.cancel') + '</button>' +
                    '<button id="saveAddUserBtn" class="btn btn-primary">' + t('common.save') + '</button>' +
                '</div>'
            );

            document.getElementById('saveAddUserBtn').addEventListener('click', async () => {
                const username = document.getElementById('newUsername').value.trim();
                const password = document.getElementById('newUserPassword').value;
                const isAdmin = document.getElementById('newUserIsAdmin').checked;

                if (!username || !password) {
                    showNotification(t('user.fillUsernamePassword'), 'error');
                    return;
                }

                try {
                    const response = await fetch('/api/users', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username, password, is_admin: isAdmin })
                    });

                    const data = await response.json();
                    if (data.success) {
                        showNotification(t('user.created'), 'success');
                        closeModal(modal);
                        loadUsersList();
                    } else {
                        showNotification(data.message || t('user.createFailed'), 'error');
                    }
                } catch (error) {
                    showNotification('Network error: ' + error.message, 'error');
                }
            });

            document.getElementById('cancelAddUserBtn').addEventListener('click', () => {
                closeModal(modal);
            });
        }

        function showEditUserModal(id, username, isAdmin) {
            const modal = createModal(t('user.editUser'), 
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: block; margin-bottom: 0.5rem; color: var(--text-primary);">' + t('connection.user') + '</label>' +
                    '<input type="text" id="editUsername" value="' + escapeHtml(username) + '" style="width: 100%; padding: 0.5rem; border: 1px solid var(--border-color); border-radius: 4px; background: var(--surface); color: var(--text-primary);">' +
                '</div>' +
                '<div style="margin-bottom: 1rem;">' +
                    '<label style="display: flex; align-items: center; gap: 0.5rem; color: var(--text-primary);">' +
                        '<input type="checkbox" id="editUserIsAdmin" ' + (isAdmin ? 'checked' : '') + '>' +
                        '<span>' + t('user.admin') + '</span>' +
                    '</label>' +
                '</div>' +
                '<div style="display: flex; gap: 0.5rem; justify-content: flex-end;">' +
                    '<button id="cancelEditUserBtn" class="btn btn-secondary">' + t('common.cancel') + '</button>' +
                    '<button id="saveEditUserBtn" class="btn btn-primary">' + t('common.save') + '</button>' +
                '</div>'
            );

            document.getElementById('saveEditUserBtn').addEventListener('click', async () => {
                const username = document.getElementById('editUsername').value.trim();
                const isAdmin = document.getElementById('editUserIsAdmin').checked;

                if (!username) {
                    showNotification(t('user.fillUsername'), 'error');
                    return;
                }

                try {
                    const response = await fetch('/api/users/' + id, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username, is_admin: isAdmin })
                    });

                    const data = await response.json();
                    if (data.success) {
                        showNotification(t('user.updated'), 'success');
                        closeModal(modal);
                        loadUsersList();
                    } else {
                        showNotification(data.message || t('user.updateFailed'), 'error');
                    }
                } catch (error) {
                    showNotification('Network error: ' + error.message, 'error');
                }
            });

            document.getElementById('cancelEditUserBtn').addEventListener('click', () => {
                closeModal(modal);
            });
        }

        function deleteUser(id) {
            fetch('/api/users/' + id, { method: 'DELETE' })
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        showNotification(t('user.deleted'), 'success');
                        loadUsersList();
                    } else {
                        showNotification(data.message || t('user.deleteFailed'), 'error');
                    }
                })
                .catch(err => {
                    showNotification('Network error: ' + err.message, 'error');
                });
        }

        function handleLogout() {
            fetch('/api/auth/logout', { method: 'POST' })
                .then(() => {
                    window.location.href = '/login';
                })
                .catch(err => {
                    console.error('Logout error:', err);
                    window.location.href = '/login';
                });
        }

        function createModal(title, content) {
            const modal = document.createElement('div');
            modal.style.cssText = 'position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 10000;';
            
            const modalContent = document.createElement('div');
            modalContent.style.cssText = 'background: var(--surface); border-radius: 8px; padding: 2rem; max-width: 500px; width: 90%; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);';
            modalContent.innerHTML = 
                '<h2 style="margin-bottom: 1.5rem; color: var(--text-primary);">' + escapeHtml(title) + '</h2>' +
                content;

            modal.appendChild(modalContent);
            document.body.appendChild(modal);

            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    closeModal(modal);
                }
            });

            return modal;
        }

        function closeModal(modal) {
            if (modal && modal.parentNode) {
                modal.parentNode.removeChild(modal);
            }
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        function showNotification(message, type) {
            if (typeof window.showNotification === 'function') {
                window.showNotification(message, type);
            } else {
                alert(message);
            }
        }
    }
})();

