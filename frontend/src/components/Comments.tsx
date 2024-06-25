import React, { useEffect, useState } from 'react';
import { Comments, getComments } from '@components/gRPC';

interface CommentProps {
    comment: Comments;
    level?: number; // Уровень вложенности комментария для отступов
}

const formatDate = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);

    const dateOptions: Intl.DateTimeFormatOptions = {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
    };

    const timeOptions: Intl.DateTimeFormatOptions = {
        hour: 'numeric',
        minute: 'numeric',
    };

    return date.toLocaleDateString(undefined, dateOptions) + ' ' + date.toLocaleTimeString(undefined, timeOptions);
}

const Comment: React.FC<CommentProps> = ({ comment, level = 0 }) => {
    const marginLeft = level * 20; // Вычисляемый отступ в зависимости от уровня вложенности

    return (
        <div className="comment" style={{marginLeft}}>
            <div className="comment__user-info">
                <img src={comment.user.photoUrl}/>
                <a href={`/user/${comment.user.username}`}>{comment.user.firstName} {comment.user.lastName}</a>
                <p>{formatDate(comment.createdAt.seconds)}</p>
            </div>
            <div>{comment.text}</div>
            {comment.children && comment.children.length > 0 && (
                <div className="comment">
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
        getComments(id, 10, 1)
            .then((commentsData) => {
                setComments(commentsData);
            })
            .catch((error) => {
                console.error('Error fetching comments:', error);
            });
    }, []);

    return (
        <>
            {comments.map((comment) => (
                <Comment key={comment.id} comment={comment} />
            ))}
        </>
    );
};

export default CommentsComponent;