import React, {memo, useState, useEffect} from 'react';
import {deleteComment} from '@components/gRPC.tsx';
import {useAuth} from '@/contexts/Auth.tsx';
import {Comments} from '@components/gRPC.tsx';
import {formatDate} from '@/utils/FormatDate.ts';

interface CommentProps {
    comment: Comments;
    level?: number;
    onReply: (parentId: number, parentComment: Comments) => void;
    onEdit: (comment: Comments) => void;
}

const Comment: React.FC<CommentProps> = memo(({comment, level = 0, onReply, onEdit}) => {
    const showReplyButton = level < 7;
    const [userCheck, setUserCheck] = useState(false);
    const {userId} = useAuth();

    useEffect(() => {
        if (comment && comment.user && userId && userId === comment.user.id) {
            setUserCheck(true);
        }
    }, [userId, comment]);

    if (!comment) {
        return null; // Если comment не существует, можно вернуть null или другую заглушку
    }

    const handleReplyClick = () => {
        onReply(comment.id, comment);
    };

    const handleEditClick = () => {
        onEdit(comment);
    };

    const handleDeleteClick = async () => {
        try {
            await deleteComment(comment.id);
            window.location.reload();
        } catch (error) {
            console.error('Failed to add comment:', error);
        }
    };

    return (
        <div className={`comment level-${level} id-${comment.id}`}>
            <div className="comment__user-info">
                <img src={comment.user.photoUrl} alt="User"/>
                <a href={`/user/${comment.user.username}`}>
                    {comment.user.firstName} {comment.user.lastName}
                </a>
                <p>{formatDate(comment.createdAt.seconds)}</p>
            </div>
            <p>{comment.text}</p>
            {showReplyButton && (
                <span className="reply-button" onClick={handleReplyClick}>
                    Ответить
                </span>
            )}
            {userCheck && (
                <>
                    <span className="edit-button" onClick={handleEditClick}>Edit</span>
                    <span className="delete-button" onClick={handleDeleteClick}>Delete</span>
                </>
            )}
            {comment.children && comment.children.length > 0 && (
                <div className="comment">
                    {comment.children.map((child) => (
                        <Comment
                            key={child.id}
                            comment={child}
                            level={level + 1}
                            onReply={onReply}
                            onEdit={onEdit}
                        />
                    ))}
                </div>
            )}
        </div>
    );
});

export default Comment;
