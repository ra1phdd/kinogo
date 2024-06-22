import React, { useEffect, useState } from 'react';
import { Comments, getComments } from '@components/gRPC';

interface CommentProps {
    comment: Comments;
    level?: number; // Уровень вложенности комментария для отступов
}

const Comment: React.FC<CommentProps> = ({ comment, level = 0 }) => {
    const marginLeft = 20 + level * 10; // Вычисляемый отступ в зависимости от уровня вложенности

    return (
        <div style={{ marginLeft, marginTop: '10px', borderLeft: '2px solid #ccc', paddingLeft: '10px' }}>
            <div>
                <strong>User ID: {comment.userId}</strong>
            </div>
            <div>{comment.text}</div>
            <div>Created At: {comment.createdAt.seconds}</div>
            {comment.children && comment.children.length > 0 && (
                <div style={{ marginTop: '10px' }}>
                    {comment.children.map((child) => (
                        <Comment key={child.id} comment={child} level={level + 1} />
                    ))}
                </div>
            )}
        </div>
    );
};

const CommentsComponent: React.FC<{ id: number }> = ({ id }) => {
    const [comments, setComments] = useState<Comments[]>([]);

    useEffect(() => {
        // Пример вызова функции getComments (предполагается, что getComments возвращает Promise<Comments[]>)
        getComments(id, 10, 1) // Пример вызова функции с конкретными параметрами
            .then((commentsData) => {
                setComments(commentsData);
            })
            .catch((error) => {
                console.error('Error fetching comments:', error);
            });
    }, []); // Пустой массив зависимостей означает, что useEffect запустится один раз при монтировании компонента

    return (
        <div>
            <h1>Comments</h1>
            {comments.map((comment) => (
                <Comment key={comment.id} comment={comment} />
            ))}
        </div>
    );
};

export default CommentsComponent;