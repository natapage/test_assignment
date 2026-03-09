<template>
  <div id="app">
    <nav class="navbar">
      <div class="navbar-brand">Goldex</div>
      <div class="navbar-links">
        <router-link to="/machines">Аппараты</router-link>
        <router-link to="/locations">Локации</router-link>
        <router-link to="/map">Карта</router-link>
        <router-link to="/statistics">Статистика</router-link>
      </div>
    </nav>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f5f5f5;
  color: #333;
}
.navbar {
  background: #1a1a2e;
  color: #fff;
  padding: 0 24px;
  height: 56px;
  display: flex;
  align-items: center;
  gap: 32px;
}
.navbar-brand {
  font-size: 20px;
  font-weight: 700;
  color: #e6b800;
}
.navbar-links {
  display: flex;
  gap: 16px;
}
.navbar-links a {
  color: #ccc;
  text-decoration: none;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background 0.2s;
}
.navbar-links a:hover,
.navbar-links a.router-link-active {
  background: rgba(255,255,255,0.1);
  color: #fff;
}
.content {
  max-width: 1200px;
  margin: 24px auto;
  padding: 0 16px;
}
table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}
th, td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}
th {
  background: #f8f9fa;
  font-weight: 600;
  font-size: 13px;
  text-transform: uppercase;
  color: #666;
}
tr:hover {
  background: #f8f9fa;
}
button {
  padding: 6px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}
.btn-primary {
  background: #1a73e8;
  color: #fff;
}
.btn-primary:hover {
  background: #1557b0;
}
.btn-danger {
  background: #dc3545;
  color: #fff;
}
.btn-danger:hover {
  background: #b02a37;
}
.btn-secondary {
  background: #6c757d;
  color: #fff;
}
.btn-secondary:hover {
  background: #565e64;
}
h1 {
  font-size: 24px;
  margin-bottom: 16px;
}
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  min-width: 400px;
  max-width: 500px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.2);
}
.modal h2 {
  margin-bottom: 16px;
}
.form-group {
  margin-bottom: 12px;
}
.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 4px;
  color: #555;
}
.form-group input,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}
.modal-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 16px;
}
.error {
  color: #dc3545;
  margin-top: 8px;
  font-size: 14px;
}
</style>
