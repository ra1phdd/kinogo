import React, { useState } from 'react';

const Editor: React.FC = () => {
    // Состояние для хранения значения textarea
    const [textValue, setTextValue] = useState<string>('');

    // Обработчик изменения значения textarea
    const handleTextareaChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        setTextValue(event.target.value);
    };

    // Обработчик нажатия на кнопку
    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        console.log('Текст из textarea:', textValue);
    };

    return (
        <form className="editor" onSubmit={handleSubmit}>
            <textarea
                value={textValue}
                onChange={handleTextareaChange}
                placeholder="Введите текст..."
            /><br/>
            <button type="submit">Отправить</button>
        </form>
    );
};

export default Editor;