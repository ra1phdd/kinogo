import React from "react";
import '@assets/styles/pages/studio/studio.css'

// Компонент Home
const Studio: React.FC = () => {
    return (
        <>
            <aside>
                <div className="aside__header">
                    <a href="/studio">Главная</a>
                    <a href="/studio/content">Контент</a>
                    <a href="/studio/comments">Комментарии</a>
                    <a href="/studio/support">Обращения</a>
                </div>
            </aside>
            <div className="sections">

            </div>
        </>
    );
};

export default Studio;
