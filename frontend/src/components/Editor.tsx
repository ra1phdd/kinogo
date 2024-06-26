import React, {memo, useEffect, useRef, useState} from 'react';
import {Comments, addComment, updateComment} from "@components/gRPC.tsx";

interface EditorProps {
    parentId: number | null;
    parentComment?: Comments | null;
    movieId: number;
    onSubmit: (text: string, parentId: number | null) => void;
    onCancelReply: () => void;
    refreshComments: () => void;
    editComment?: Comments | null;
}

const Editor: React.FC<EditorProps> = memo(({ parentId, parentComment, movieId, onSubmit, onCancelReply, refreshComments, editComment }) => {
    const [textValue, setTextValue] = useState<string>('');
    const formRef = useRef<HTMLFormElement>(null);

    const handleTextareaChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        setTextValue(event.target.value);
    };

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            onSubmit(textValue, parentId);
            let userId: number | null = null;

            const storedUserId = localStorage.getItem('userid');
            if (storedUserId !== null) {
                userId = Number(storedUserId);
            } else {
                console.error('User ID is not available.');
                return;
            }

            if (editComment) {
                await updateComment(editComment.id, textValue);
            } else {
                await addComment(parentId, movieId, userId, textValue);
            }
            setTextValue('');

            refreshComments();
        } catch (error) {
            console.error('Failed to add comment:', error);
        }
    };

    useEffect(() => {
        if (parentId !== null && formRef.current) {
            formRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }, [parentId]);

    useEffect(() => {
        if (editComment) {
            setTextValue(editComment.text);
            if (formRef.current) {
                formRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    }, [editComment]);

    return (
        <form ref={formRef} className="editor" onSubmit={handleSubmit}>
            {parentComment && (
                <div className="editor__comment">
                    <p className="editor__comment-answer"><span className="editor__comment-cancel" onClick={onCancelReply}>×</span> Ответ
                        пользователю <img src={parentComment.user.photoUrl} alt="User"/>{parentComment.user.firstName} {parentComment.user.lastName}</p>
                    <p className="editor__comment-text">{parentComment.text}</p>
                </div>
            )}
            {editComment && (
                <div className="editor__comment">
                    <p className="editor__comment-answer"><span className="editor__comment-cancel" onClick={onCancelReply}>×</span> Изменение
                        текста комментария</p>
                    <p className="editor__comment-text">{editComment.text}</p>
                </div>
            )}
            <textarea
                value={textValue}
                onChange={handleTextareaChange}
                placeholder="Введите ваш комментарий..."
            /><br />
            <button type="submit">Отправить</button>
        </form>
    );
});

export default Editor;