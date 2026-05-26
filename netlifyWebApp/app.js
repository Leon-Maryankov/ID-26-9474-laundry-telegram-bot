// Telegram Web App initialization
const tg = window.Telegram.WebApp;
tg.expand();

// API Configuration
const API_TOKEN = 'token';
const API_BASE_URL = 'http://localhost:8080';

// Application State
let appState = {
    user: null,
    isRegistered: false,
    currentTab: 'available',
    slots: {
        available: [],
        booked: []
    }
};

// DOM Elements
const screens = {
    loading: document.getElementById('loadingScreen'),
    registration: document.getElementById('registrationScreen'),
    main: document.getElementById('mainScreen')
};

const forms = {
    registration: document.getElementById('registrationForm')
};

const lists = {
    available: document.getElementById('availableSlotsList'),
    booked: document.getElementById('bookedSlotsList')
};

const modals = {
    profile: document.getElementById('profileModal')
};

const toasts = {
    error: document.getElementById('errorToast'),
    success: document.getElementById('successToast')
};

// Utility Functions
function showScreen(screenName) {
    Object.values(screens).forEach(screen => screen.classList.remove('active'));
    screens[screenName].classList.add('active');
}

function showToast(type, message) {
    const toast = toasts[type];
    const messageElement = document.getElementById(`${type}ToastMessage`);
    messageElement.textContent = message;
    toast.classList.add('visible');
    setTimeout(() => {
        toast.classList.remove('visible');
    }, 3000);
}

// Универсальная функция форматирования даты
function formatDateTime(timestamp) {
    // Если timestamp – число (Unix timestamp в секундах)
    let ts = typeof timestamp === 'number' ? timestamp : Number(timestamp);
    if (!isNaN(ts) && ts > 0) {
        const date = new Date(ts * 1000);
        if (!isNaN(date.getTime())) {
            return formatDateFromDate(date);
        }
    }
    // Если timestamp – строка ISO
    const date = new Date(timestamp);
    if (!isNaN(date.getTime())) {
        return formatDateFromDate(date);
    }
    console.error('Invalid timestamp:', timestamp);
    return 'Неверная дата';
}

function formatDateFromDate(date) {
    const days = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];
    const months = ['янв', 'фев', 'мар', 'апр', 'май', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек'];
    const dayOfWeek = days[date.getDay()];
    const day = date.getDate();
    const month = months[date.getMonth()];
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${dayOfWeek}, ${day} ${month} ${hours}:${minutes}`;
}

function getMachineTypeLabel(type) {
    return type === 'washing' ? '🧺 Стиральная машина' : '🌡️ Сушилка';
}

function getMachineTypeColor(type) {
    return type === 'washing' ? '#007AFF' : '#FF9500';
}

// API Functions
async function apiCall(endpoint, method = 'GET', data = null) {
    const url = `${API_BASE_URL}/${API_TOKEN}${endpoint}`;
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json'
        }
    };
    if (data) {
        options.body = JSON.stringify(data);
    }
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        showToast('error', 'Ошибка подключения к серверу');
        throw error;
    }
}

async function checkUserExists(telegramId) {
    try {
        const response = await apiCall('/doesUserExist', 'POST', { telegramId: telegramId });
        if (response.status === 'good' && response.answer) {
            return response.user;
        }
        return null;
    } catch (error) {
        console.error('Error checking user:', error);
        return null;
    }
}

async function registerUser(name, surname, email) {
    try {
        const telegramId = tg.initDataUnsafe.user?.id;
        if (!telegramId) {
            showToast('error', 'Не удалось определить Telegram ID');
            return null;
        }
        const response = await apiCall('/registerUser', 'POST', {
            telegramId: telegramId,
            name: name,
            surname: surname,
            email: email
        });
        if (response.status === 'good') {
            showToast('success', 'Успешная регистрация!');
            return { id: telegramId, name, surname, email };
        } else {
            showToast('error', 'Ошибка регистрации');
            return null;
        }
    } catch (error) {
        console.error('Error registering user:', error);
        showToast('error', 'Ошибка при регистрации');
        return null;
    }
}

async function getAvailableSlots() {
    try {
        const response = await apiCall('/getAvailableSlots', 'POST', {});
        if (response.status === 'good') {
            return response.Slots || [];
        }
        return [];
    } catch (error) {
        console.error('Error fetching available slots:', error);
        return [];
    }
}

async function getBookedSlots(userId) {
    try {
        const response = await apiCall('/getSlotsAttachedToUser', 'POST', { userId: userId });
        if (response.status === 'good') {
            return response.slots || [];
        }
        return [];
    } catch (error) {
        console.error('Error fetching booked slots:', error);
        return [];
    }
}

async function bookSlot(slotId) {
    if (!appState.user) return false;
    try {
        const response = await apiCall('/occupySlot', 'POST', {
            slotId: slotId,
            userId: appState.user.id
        });
        if (response.status === 'good') {
            showToast('success', 'Слот успешно забронирован!');
            return true;
        } else {
            const reasonMap = {
                'too many washing machines': 'Вы можете забронировать максимум 2 стиральные машины',
                'too many dryer machines': 'Вы можете забронировать максимум 1 сушилку'
            };
            const errorMessage = reasonMap[response.reason] || 'Не удалось забронировать слот';
            showToast('error', errorMessage);
            return false;
        }
    } catch (error) {
        console.error('Error booking slot:', error);
        showToast('error', 'Ошибка при бронировании');
        return false;
    }
}

async function cancelSlot(slotId) {
    if (!appState.user) return false;
    try {
        const response = await apiCall('/cancelSlot', 'POST', {
            slotId: slotId,
            userId: appState.user.id
        });
        if (response.status === 'good') {
            showToast('success', 'Бронирование отменено');
            return true;
        } else {
            showToast('error', 'Не удалось отменить бронирование');
            return false;
        }
    } catch (error) {
        console.error('Error canceling slot:', error);
        showToast('error', 'Ошибка при отмене бронирования');
        return false;
    }
}

// UI Rendering Functions
function renderAvailableSlots(slots) {
    const machineTypeFilter = document.getElementById('machineTypeFilter').value;
    const dateFilter = document.getElementById('dateFilter').value;
    let filteredSlots = slots;
    if (machineTypeFilter) {
        filteredSlots = filteredSlots.filter(s => s.MachineType === machineTypeFilter);
    }
    if (dateFilter) {
        const filterDate = new Date(dateFilter);
        filteredSlots = filteredSlots.filter(s => {
            const slotDate = new Date(s.Date * 1000);
            return slotDate.toDateString() === filterDate.toDateString();
        });
    }
    if (filteredSlots.length === 0) {
        lists.available.innerHTML = '<div class="empty-message"><p>😴 Нет доступных слотов с такими фильтрами</p></div>';
        return;
    }
    const groupedByDate = {};
    filteredSlots.forEach(slot => {
        const dateStr = formatDateTime(slot.Date).split(' ').slice(0, 3).join(' ');
        if (!groupedByDate[dateStr]) {
            groupedByDate[dateStr] = [];
        }
        groupedByDate[dateStr].push(slot);
    });
    let html = '';
    Object.entries(groupedByDate).forEach(([dateStr, daySlots]) => {
        html += `<div class="slots-group">
            <h3 class="slots-group-title">${dateStr}</h3>
            <div class="slots-grid">`;
        daySlots.forEach(slot => {
            const time = formatDateTime(slot.Date);
            const machineLabel = getMachineTypeLabel(slot.MachineType);
            const color = getMachineTypeColor(slot.MachineType);
            html += `<div class="slot-card available">
                <div class="slot-time">${time}</div>
                <div class="slot-type" style="border-left: 4px solid ${color}">
                    ${machineLabel}
                </div>
                <div class="slot-machine">Машина #${slot.MachineId + 1}</div>
                <button class="button button-primary slot-book-btn" data-slot-id="${slot.Id}">
                    Забронировать
                </button>
            </div>`;
        });
        html += '</div></div>';
    });
    lists.available.innerHTML = html;
    document.querySelectorAll('.slot-book-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const slotId = parseInt(e.target.dataset.slotId);
            btn.disabled = true;
            btn.textContent = 'Обработка...';
            const success = await bookSlot(slotId);
            if (success) {
                await loadAllData();
            }
            btn.disabled = false;
            btn.textContent = 'Забронировать';
        });
    });
}

function renderBookedSlots(slots) {
    if (slots.length === 0) {
        lists.booked.innerHTML = '<div class="empty-message"><p>📭 У вас нет активных бронирований</p></div>';
        return;
    }
    const sortedSlots = slots.sort((a, b) => a.Date - b.Date);
    let html = '';
    sortedSlots.forEach(slot => {
        const time = formatDateTime(slot.Date);
        const machineLabel = getMachineTypeLabel(slot.MachineType);
        const color = getMachineTypeColor(slot.MachineType);
        const now = Date.now() / 1000;
        const isPast = slot.Date < now;
        html += `<div class="slot-card booked ${isPast ? 'past' : ''}">
            <div class="slot-time">${time}</div>
            <div class="slot-type" style="border-left: 4px solid ${color}">
                ${machineLabel}
            </div>
            <div class="slot-machine">Машина #${slot.MachineId + 1}</div>
            ${!isPast ? `<button class="button button-danger slot-cancel-btn" data-slot-id="${slot.Id}">
                Отменить
            </button>` : '<div class="slot-status-past">Завершено</div>'}
        </div>`;
    });
    lists.booked.innerHTML = html;
    document.querySelectorAll('.slot-cancel-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const slotId = parseInt(e.target.dataset.slotId);
            btn.disabled = true;
            btn.textContent = 'Отмена...';
            const success = await cancelSlot(slotId);
            if (success) {
                await loadAllData();
            }
            btn.disabled = false;
            btn.textContent = 'Отменить';
        });
    });
}

// Data Loading
async function loadAllData() {
    try {
        appState.slots.available = await getAvailableSlots();
        appState.slots.booked = await getBookedSlots(appState.user.id);
        renderAvailableSlots(appState.slots.available);
        renderBookedSlots(appState.slots.booked);
    } catch (error) {
        console.error('Error loading data:', error);
    }
}

// Event Listeners Setup
function setupEventListeners() {
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const tab = e.currentTarget.dataset.tab;
            appState.currentTab = tab;
            document.querySelectorAll('.tab-button').forEach(b => b.classList.remove('tab-active'));
            e.currentTarget.classList.add('tab-active');
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('tab-active'));
            document.getElementById(`${tab}Tab`).classList.add('tab-active');
        });
    });
    document.getElementById('machineTypeFilter').addEventListener('change', () => {
        renderAvailableSlots(appState.slots.available);
    });
    document.getElementById('dateFilter').addEventListener('change', () => {
        renderAvailableSlots(appState.slots.available);
    });
    document.getElementById('profileBtn').addEventListener('click', () => {
        document.getElementById('profileName').textContent = appState.user.name;
        document.getElementById('profileSurname').textContent = appState.user.surname;
        document.getElementById('profileEmail').textContent = appState.user.email;
        modals.profile.classList.add('visible');
    });
    document.querySelector('.modal-close').addEventListener('click', () => {
        modals.profile.classList.remove('visible');
    });
    modals.profile.addEventListener('click', (e) => {
        if (e.target === modals.profile) {
            modals.profile.classList.remove('visible');
        }
    });
    forms.registration.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('nameInput').value.trim();
        const surname = document.getElementById('surnameInput').value.trim();
        const email = document.getElementById('emailInput').value.trim();
        if (!name || !surname || !email) {
            showToast('error', 'Пожалуйста, заполните все поля');
            return;
        }
        const submitBtn = forms.registration.querySelector('button[type="submit"]');
        submitBtn.disabled = true;
        submitBtn.textContent = 'Регистрация...';
        const user = await registerUser(name, surname, email);
        if (user) {
            appState.user = user;
            appState.isRegistered = true;
            document.getElementById('userName').textContent = `${user.name} ${user.surname}`;
            showScreen('main');
            await loadAllData();
        }
        submitBtn.disabled = false;
        submitBtn.textContent = 'Зарегистрироваться';
    });
    document.querySelectorAll('.toast-close').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.target.closest('.toast').classList.remove('visible');
        });
    });
}

// Initialize Application
async function initializeApp() {
    try {
        showScreen('loading');
        const telegramId = tg.initDataUnsafe.user?.id;
        if (!telegramId) {
            showToast('error', 'Не удалось получить данные пользователя из Telegram');
            showScreen('registration');
            return;
        }
        const existingUser = await checkUserExists(telegramId);
        if (existingUser) {
            appState.user = existingUser;
            appState.isRegistered = true;
            document.getElementById('userName').textContent = `${existingUser.name} ${existingUser.surname}`;
            showScreen('main');
            setupEventListeners();
            await loadAllData();
        } else {
            showScreen('registration');
            setupEventListeners();
        }
    } catch (error) {
        console.error('Initialization error:', error);
        showToast('error', 'Ошибка инициализации приложения');
        showScreen('registration');
        setupEventListeners();
    }
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeApp);
} else {
    initializeApp();
}