import React, { memo, useEffect } from 'react';
import { Comments, addComment, updateComment } from "@components/gRPC.tsx";
import { useAuth } from "@/hooks/useAuth.ts";
import { useTextareaState } from '@/hooks/useTextareaState';
import { useScrollIntoView } from '@/hooks/useScrollIntoView';

interface EditorProps {
    parentId: number | null;
    parentComment?: Comments | null;
    movieId: number;
    onSubmit: (text: string, parentId: number | null) => void;
    onCancelReply: () => void;
    editComment?: Comments | null;
}

const Editor: React.FC<EditorProps> = memo(({ parentId, parentComment, movieId, onSubmit, onCancelReply, editComment }) => {
    const { value: textValue, onChange: handleTextareaChange, reset: resetTextarea } = useTextareaState(editComment?.text);
    const { userId } = useAuth();
    const formRef = useScrollIntoView(parentId !== null || (editComment !== null && editComment !== undefined));

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            onSubmit(textValue, parentId);

            if (!editComment && userId !== null && textValue.trim() !== '') {
                await addComment(parentId, movieId, userId, textValue);
            } else if (editComment) {
                await updateComment(editComment.id, textValue);
            }

            resetTextarea();
            window.location.reload();
        } catch (error) {
            console.error('Failed to add/update comment:', error);
        }
    };

    useEffect(() => {
        resetTextarea();
        if (editComment) {
            formRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
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
