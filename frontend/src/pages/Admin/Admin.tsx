import React from "react";
import '@assets/styles/pages/admin/home.css'

// Компонент Home
const Admin: React.FC = () => {
    return (
        <>
            <aside>
                <div className="aside__header">
                    <a href="/admin">Главная</a>
                    <a href="/admin/content">Контент</a>
                    <a href="/admin/stats">Статистика</a>
                    <a href="/admin/comments">Комментарии</a>
                    <a href="/admin/support">Обращения</a>
                </div>
            </aside>
            <div className="sections">

            </div>
        </>
    );
};

export default Admin;
