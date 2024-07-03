import React, { memo, useEffect, useRef, useState } from 'react';
import { Comments, addComment, updateComment } from "@components/gRPC.tsx";
import { GetUserId } from "@components/JwtDecode.tsx";

interface EditorProps {
    parentId: number | null;
    parentComment?: Comments | null;
    movieId: number;
    onSubmit: (text: string, parentId: number | null) => void;
    onCancelReply: () => void;
    editComment?: Comments | null;
}

const Editor: React.FC<EditorProps> = memo(({ parentId, parentComment, movieId, onSubmit, onCancelReply, editComment }) => {
    const [textValue, setTextValue] = useState<string>('');
    const formRef = useRef<HTMLFormElement>(null);
    const userId = GetUserId();

    const handleTextareaChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        setTextValue(event.target.value);
    };

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            onSubmit(textValue, parentId);

            if (!editComment && userId !== undefined && textValue.trim() !== '') {
                await addComment(parentId, movieId, userId, textValue);
            } else if (editComment) {
                await updateComment(editComment.id, textValue);
            }

            setTextValue('');
        } catch (error) {
            console.error('Failed to add/update comment:', error);
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
                    <p className="editor__comment-answer">
                        <span className="editor__comment-cancel" onClick={onCancelReply}>×</span> Ответ пользователю
                        <img src={parentComment.user.photoUrl} alt="User"/> {parentComment.user.firstName} {parentComment.user.lastName}
                    </p>
                    <p className="editor__comment-text">{parentComment.text}</p>
                </div>
            )}
            {editComment && (
                <div className="editor__comment">
                    <p className="editor__comment-answer">
                        <span className="editor__comment-cancel" onClick={onCancelReply}>×</span> Изменение текста комментария
                    </p>
                </div>
            )}
            <textarea
                value={textValue}
                onChange={handleTextareaChange}
                placeholder="Введите ваш комментарий..."
            />
            <br />
            <button type="submit">
                <div id="circle"></div>
                <span>Отправить</span>
            </button>
        </form>
    );
});

export default Editor;
