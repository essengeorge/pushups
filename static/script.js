// bebebe
const API_URL = 'https://pushups.gveserver.ru';
let chart = null;
let allStats = [];
let token = localStorage.getItem('token');
let currentUsername = localStorage.getItem('username');

// ============= AUTH =============
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (response.ok && data.token) {
            localStorage.setItem('token', data.token);
            localStorage.setItem('username', username);
            token = data.token;
            currentUsername = username;
            showApp();
        } else {
            document.getElementById('loginError').textContent = 'Неверные данные';
        }
    } catch (error) {
        document.getElementById('loginError').textContent = 'Ошибка сети';
    }
});

document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(`${API_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (response.ok && data.status === 'ok') {
            document.getElementById('registerError').innerHTML =
                '<div class="success">Аккаунт создан! Теперь можно войти</div>';
            setTimeout(() => {
                toggleForms();
                document.getElementById('registerForm').reset();
            }, 1500);
        } else {
            document.getElementById('registerError').textContent = 'Ошибка регистрации';
        }
    } catch (error) {
        document.getElementById('registerError').textContent = 'Ошибка сети';
    }
});

document.getElementById('toggleBtn').addEventListener('click', toggleForms);

function toggleForms() {
    document.getElementById('loginForm').classList.toggle('active');
    document.getElementById('registerForm').classList.toggle('active');
    const toggleText = document.getElementById('toggleText');
    const toggleBtn = document.getElementById('toggleBtn');

    if (document.getElementById('loginForm').classList.contains('active')) {
        toggleText.textContent = 'Нет аккаунта? ';
        toggleBtn.textContent = 'Создать';
    } else {
        toggleText.textContent = 'Уже есть аккаунт? ';
        toggleBtn.textContent = 'Войти';
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    token = null;
    currentUsername = null;
    showAuth();
}


// ============= APP =============
async function showApp() {
    document.querySelector('.auth-page').classList.remove('active');
    document.querySelector('.app-page').classList.add('active');
    document.getElementById('currentUser').textContent = currentUsername;
    await loadStats();
    await loadFriends();
}

function showAuth() {
    document.querySelector('.auth-page').classList.add('active');
    document.querySelector('.app-page').classList.remove('active');
    document.getElementById('loginForm').reset();
    document.getElementById('registerForm').reset();
    document.getElementById('loginForm').classList.add('active');
    document.getElementById('registerForm').classList.remove('active');
}

async function addPushups() {
    if (!token) return;

    const count = parseInt(document.getElementById('pushupsCount').value);
    if (!count || count < 1) return;

    try {
        const response = await fetch(`${API_URL}/pushups`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ count })
        });

        if (response.ok) {
            showMessage('✅ Отжимания добавлены!', 'success');
            document.getElementById('pushupsCount').value = '10';
            await loadStats();
        } else {
            showMessage('❌ Ошибка при добавлении', 'error');
        }
    } catch (error) {
        showMessage('❌ Ошибка сети', 'error');
    }
}

async function loadStats() {
    if (!token) return;

    try {
        const response = await fetch(`${API_URL}/stats`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.ok) {
            allStats = await response.json();
            updateChart();
            updateLeaderboard();
        } else {
            showMessage('❌ Ошибка загрузки статистики', 'error');
        }
    } catch (error) {
        showMessage('❌ Ошибка сети', 'error');
    }
}

function updateChart() {
    if (!allStats || allStats.length === 0) {
        console.warn('Нет данных для графика');
        return;
    }

    const ctx = document.getElementById('statsChart');
    if (!ctx) {
        console.error('Элемент #statsChart не найден');
        return;
    }

    const colors = generateColors(allStats.length);

    // Собираем все уникальные даты
    const allDates = new Set();
    allStats.forEach(stat => {
        stat.points.forEach(point => {
            allDates.add(point.date);
        });
    });
    const sortedDates = Array.from(allDates).sort();

    // Создаём датасеты
    const datasets = allStats.map((userStats, index) => {
        const color = userStats.color || colors[index];
        const dateMap = {};
        userStats.points.forEach(point => {
            dateMap[point.date] = point.value;
        });
        return {
            label: userStats.username,
            data: sortedDates.map(date => dateMap[date] || 0), // for line null, for bar 0
            borderColor: color,
            backgroundColor: color + '80', // for bar '80', for line '20'
            borderWidth: 2,
            borderSkipped: false,
/*
// for line uncomment
            tension: 0.4,
            fill: true,
            pointRadius: 4,
            pointBackgroundColor: color,
            pointBorderColor: '#fff',
            pointBorderWidth: 2,
            pointHoverRadius: 6,
            spanGaps: true
// for line uncomment
*/
        };
    });

    if (chart) {
        chart.destroy();
    }

    chart = new Chart(ctx, {
        type: 'bar', // 'bar'/'line'
        data: {
            labels: sortedDates,
            datasets: datasets
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                mode: 'index',
                intersect: false
            },
            plugins: {
                legend: {
                    position: 'top',
                    labels: {
                        font: { size: 13 },
                        padding: 15,
                        usePointStyle: true
                    }
                },
                tooltip: {
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    padding: 12,
                    titleFont: { size: 14 },
                    bodyFont: { size: 13 },
                    borderColor: '#ddd',
                    borderWidth: 1
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grace: '5%',
                    grid: {
                        color: '#f0f0f0',
                        drawBorder: false
                    },
                    ticks: {
                        font: { size: 12 }
                    },
                    title: {
                        display: true,
                        text: 'Отжимания',
                        font: { size: 13, weight: 'bold' }
                    }
                },
                x: {
                    grid: {
                        display: false,
                        drawBorder: false
                    },
                    ticks: {
                        font: { size: 12 }
                    }
                }
            }
        }
    });
}

function updateLeaderboard() {
    const leaderboard = document.getElementById('leaderboard');
    leaderboard.innerHTML = '';

    // Считаем общее количество отжиманий для каждого пользователя
    const totals = allStats.map(stat => ({
        username: stat.username,
        color: stat.color,
        total: stat.points.reduce((acc, point) => acc + point.value, 0)
    })).sort((a, b) => b.total - a.total);

    totals.forEach(item => {
        const div = document.createElement('div');
        div.className = 'leaderboard-item';
        div.innerHTML = `
            <div class="color-dot" style="background-color: ${item.color}"></div>
            <div class="leaderboard-name">${item.username}</div>
            <div class="leaderboard-count">${item.total}</div>
        `;
        leaderboard.appendChild(div);
    });
}

function showMessage(text, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = text;
    messageDiv.style.display = 'block';

    setTimeout(() => {
        messageDiv.style.display = 'none';
    }, 3000);
}

function generateColors(count) {
    const hues = [];
    for (let i = 0; i < count; i++) {
        hues.push((i * 360 / count) % 360);
    }
    return hues.map(hue => `hsl(${hue}, 70%, 50%)`);
}

function makeBtn(text, fn) {
    const btn = document.createElement('button');
    btn.className = 'add-btn';
    btn.innerText = text;
    btn.onclick = fn;
    return btn;
}

async function loadFriends() {
    if (!token) return;

    try {
        // Загружаем список друзей
        const friendsRes = await fetch(`${API_URL}/friends/list`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (friendsRes.ok) {
            const friends = await friendsRes.json();
            document.getElementById('friendsFriends').innerHTML = '<div class="panel-title">🤝Friends</div>';
            if (friends && friends.length > 0) {
                friends.forEach(f => addFriendDiv(f.username));
            }
        }

        // Входящие запросы
        const incomingRes = await fetch(`${API_URL}/friends/incoming`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (incomingRes.ok) {
            const incoming = await incomingRes.json();
            document.getElementById('friendsIncoming').innerHTML = '<div class="panel-title">📥Incoming</div>';
            if (incoming && incoming.length > 0) {
                incoming.forEach(req => addIncomingDiv(req.username));
            }
        }

        // Исходящие запросы
        const outgoingRes = await fetch(`${API_URL}/friends/outgoing`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (outgoingRes.ok) {
            const outgoing = await outgoingRes.json();
            document.getElementById('friendsPending').innerHTML = '<div class="panel-title">⏳Pending</div>';
            if (outgoing && outgoing.length > 0) {
                outgoing.forEach(req => addPendingDiv(req.username));
            }
        }

        // Заблокированные
        const blockedRes = await fetch(`${API_URL}/friends/blocked`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (blockedRes.ok) {
            const blocked = await blockedRes.json();
            document.getElementById('friendsBlocked').innerHTML = '<div class="panel-title">🚫Blocked</div>';
            if (blocked && blocked.length > 0) {
                blocked.forEach(u => addBlockedDiv(u.username));
            }
        }
    } catch (error) {
        console.error('Ошибка загрузки друзей:', error);
    }
}

async function removeFriend(username) {
    if (!token) return;
    try {
        const response = await fetch(`${API_URL}/friends/remove`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username: username })
        });
        if (response.ok) {
            showMessage('✅ Друг удален', 'success');
            document.getElementById(username).remove();
        }
    } catch (error) {
        showMessage('❌ Ошибка', 'error');
    }
}

async function blockUser(username) {
    if (!token) return;
    try {
        const response = await fetch(`${API_URL}/friends/block`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username: username })
        });
        if (response.ok) {
            showMessage('✅ Пользователь заблокирован', 'success');
            document.getElementById(username).remove();
            await loadFriends();
        }
    } catch (error) {
        showMessage('❌ Ошибка', 'error');
    }
}

async function acceptRequest(username) {
    if (!token) return;
    try {
        const response = await fetch(`${API_URL}/friends/accept`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username: username })
        });
        if (response.ok) {
            showMessage('✅ Запрос принят', 'success');
            document.getElementById(username).remove();
            await loadFriends();
        }
    } catch (error) {
        showMessage('❌ Ошибка', 'error');
    }
}

async function rejectRequest(username) {
    if (!token) return;
    try {
        const response = await fetch(`${API_URL}/friends/reject`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username: username })
        });
        if (response.ok) {
            showMessage('✅ Запрос отклонен', 'success');
            document.getElementById(username).remove();
        }
    } catch (error) {
        showMessage('❌ Ошибка', 'error');
    }
}

async function removeRequest(username) {
    if (!token) return;
    try {
        const response = await fetch(`${API_URL}/friends/reject`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username: username })
        });
        if (response.ok) {
            showMessage('✅ Запрос отменен', 'success');
            document.getElementById(username).remove();
        }
    } catch (error) {
        showMessage('❌ Ошибка', 'error');
    }
}

async function unblockUser(username) {
    if (!token) return;
    try {
        const response = await fetch(`${API_URL}/friends/unblock`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username: username })
        });
        if (response.ok) {
            showMessage('✅ Пользователь разблокирован', 'success');
            document.getElementById(username).remove();
            await loadFriends();
        }
    } catch (error) {
        showMessage('❌ Ошибка', 'error');
    }
}


// Заменить последние 4 строки на:
// (они удалят тестовые друзей и загрузят реальных)

async function sendFriendRequest() {
    if (!token) return;

    const target_username = document.getElementById('friendName').value;
    if (!target_username) return;

    try {
        const response = await fetch(`${API_URL}/friends/send`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username })
        });

        if (response.ok) {
            showMessage('✅ Запрос отправлен!', 'success');
            document.getElementById('friendName').value = '';
            await loadStats();
        } else {
            showMessage('❌ Ошибка при добавлении', 'error');
        }
    } catch (error) {
        showMessage('❌ Ошибка сети', 'error');
    }
}

async function acceptRequest(target_username) {
    if (!token) return;

    if (!target_username) return;

    try {
        const response = await fetch(`${API_URL}/friends/accept`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_username })
        });

        if (response.ok) {
            showMessage('✅ Запрос отправлен!', 'success');
            document.getElementById('friendName').value = '';
            await loadStats();
        } else {
            showMessage('❌ Ошибка при добавлении', 'error');
        }
    } catch (error) {
        showMessage('❌ Ошибка сети', 'error');
    }
}

function addFriendDiv(username) {
    const div = document.createElement('div');
    div.className = 'leaderboard-item';
    div.id = username;
    const nameDiv = document.createElement('div');
    nameDiv.className = 'leaderboard-name';
    nameDiv.textContent = username;
    div.appendChild(nameDiv);
    const removeBtn = makeBtn('Remove', () => {
        removeFriend(username);
    });
    div.appendChild(removeBtn);
    const blockBtn = makeBtn('Block', () => {
        blockUser(username);
    });
    div.appendChild(blockBtn);
    div.style.marginTop = '15px';
    document.getElementById('friendsFriends').appendChild(div);
}

function addIncomingDiv(username) {
    const div = document.createElement('div');
    div.className = 'leaderboard-item';
    div.id = username;
    const nameDiv = document.createElement('div');
    nameDiv.className = 'leaderboard-name';
    nameDiv.textContent = username;
    div.appendChild(nameDiv);
    const acceptBtn = makeBtn('Accept', () => {
        acceptRequest(username);
    });
    div.appendChild(acceptBtn);
    const rejectBtn = makeBtn('Reject', () => {
        rejectRequest(username);
    });
    div.appendChild(rejectBtn);
    const blockBtn = makeBtn('Block', () => {
        blockUser(username);
    });
    div.appendChild(blockBtn);
    div.style.marginTop = '15px';
    document.getElementById('friendsIncoming').appendChild(div);
}

function addPendingDiv(username) {
    const div = document.createElement('div');
    div.className = 'leaderboard-item';
    div.id = username;
    const nameDiv = document.createElement('div');
    nameDiv.className = 'leaderboard-name';
    nameDiv.textContent = username;
    div.appendChild(nameDiv);
    const removeReqBtn = makeBtn('Remove request', () => {
        removeRequest(username);
    });
    div.appendChild(removeReqBtn);
    div.style.marginTop = '15px';
    document.getElementById('friendsPending').appendChild(div);
}

function addBlockedDiv(username) {
    const div = document.createElement('div');
    div.className = 'leaderboard-item';
    div.id = username;
    const nameDiv = document.createElement('div');
    nameDiv.className = 'leaderboard-name';
    nameDiv.textContent = username;
    div.appendChild(nameDiv);
    const unblockBtn = makeBtn('Unblock', () => {
        unblockUser(username);
    });
    div.appendChild(unblockBtn);
    div.style.marginTop = '15px';
    document.getElementById('friendsBlocked').appendChild(div);
}

// ============= INIT =============
if (token && currentUsername) {
    showApp();
} else {
    showAuth();
}

// Обновляем статистику каждые 60 секунд
setInterval(() => {
    if (token) {
        loadStats();
    }
}, 1 * 60 * 1000);

addFriendDiv('Лёша');
addIncomingDiv('Коля');
addPendingDiv('Ваня');
addBlockedDiv('Дима');
