import React, { memo, useCallback, useEffect, useState } from 'react';
import {Comments, getComments, deleteComment} from '@components/gRPC';
import Editor from '@components/Editor.tsx';

interface CommentProps {
    comment: Comments;
    level?: number;
    onReply: (parentId: number, parentComment: Comments) => void;
    onEdit: (comment: Comments) => void;
    refreshComments: () => void;
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

    return (
        date.toLocaleDateString(undefined, dateOptions) +
        ' ' +
        date.toLocaleTimeString(undefined, timeOptions)
    );
};

const CommentsComponent: React.FC<{ id: number }> = ({ id }) => {
    const [comments, setComments] = useState<Comments[]>([]);
    const [replyToParentId, setReplyToParentId] = useState<number | null>(null);
    const [replyToParentComment, setReplyToParentComment] = useState<Comments | null>(null);
    const [editComment, setEditComment] = useState<Comments | null>(null);
    const [hasToken, setHasToken] = useState(false);
    const [movieId, setMovieId] = useState<number>(0);
    const [loading, setLoading] = useState(false);
    const [end, setEnd] = useState(false);
    const [page, setPage] = useState(1);

    useEffect(() => {
        const token = localStorage.getItem('token');
        setHasToken(!!token);
        setMovieId(id);
    }, [id]);

    const handleReply = useCallback((parentId: number, parentComment: Comments) => {
        setReplyToParentId(parentId);
        setReplyToParentComment(parentComment);
        setEditComment(null);
    }, []);

    const handleEdit = useCallback((comment: Comments) => {
        setEditComment(comment);
        setReplyToParentId(null);
        setReplyToParentComment(null);
    }, []);

    const handleCancelReply = useCallback(() => {
        setReplyToParentId(null);
        setReplyToParentComment(null);
        setEditComment(null);
    }, []);

    const handleSubmit = useCallback(() => {
        handleCancelReply();
    }, [handleCancelReply]);

    const loadMoreComments = useCallback(() => {
        if (!loading && !end) {
            setLoading(true);
            getComments(id, 10, page + 1)
                .then((newComments) => {
                    setComments(prevComments => [...prevComments, ...newComments]);
                    setPage((prevPage) => prevPage + 1);
                })
                .catch((error) => {
                    if (error.code === 5) {
                        setEnd(true);
                        return
                    }
                    console.error('Error fetching comments:', error);
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [id, loading, page]);

    const handleScroll = useCallback(() => {
        if (window.innerHeight + document.documentElement.scrollTop >= document.documentElement.offsetHeight - 100) {
            loadMoreComments();
        }
    }, [loadMoreComments]);

    useEffect(() => {
        window.addEventListener('scroll', handleScroll);
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [handleScroll]);

    const refreshComments = useCallback(() => {
        getComments(id, 10, 1)
            .then((commentsData) => {
                setComments(commentsData);
            })
            .catch((error) => {
                console.error('Error fetching comments:', error);
            });
    }, [id]);

    useEffect(() => {
        refreshComments();
    }, [id, refreshComments]);

    return (
        <>
            {hasToken && (
                <Editor
                    parentId={replyToParentId}
                    parentComment={replyToParentComment}
                    movieId={movieId}
                    onSubmit={handleSubmit}
                    onCancelReply={handleCancelReply}
                    refreshComments={refreshComments}
                    editComment={editComment}
                />
            )}
            {comments.map((comment) => (
                <Comment key={comment.id} comment={comment} onReply={handleReply} onEdit={handleEdit} refreshComments={refreshComments} />
            ))}
            {loading && <p>Загрузка...</p>}
            {!loading && comments.length === 0 && <p>No comments found.</p>}
        </>
    );
};

const Comment: React.FC<CommentProps> = memo(({ comment, level = 0, onReply, onEdit, refreshComments }) => {
    const showReplyButton = level < 7;
    const [userCheck, setUserCheck] = useState(false);
    let userId: number | null = null;

    const handleReplyClick = () => {
        onReply(comment.id, comment);
    };

    const handleEditClick = () => {
        onEdit(comment);
    };

    const handleDeleteClick = async () => {
        try {
            await deleteComment(comment.id);
            refreshComments();
        } catch (error) {
            console.error('Failed to add comment:', error);
        }
    };

    useEffect(() => {
        const storedUserId = localStorage.getItem('userid');
        if (storedUserId !== null) {
            userId = Number(storedUserId);
        } else {
            console.error('User ID is not available.');
            return;
        }

        if (userId == comment.user.id){
            setUserCheck(true);
        }
    }, []);

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
                <span className="edit-button" onClick={handleEditClick}>Edit</span>
            )}
            {userCheck && (
                <span className="delete-button" onClick={handleDeleteClick}>Delete</span>
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
                            refreshComments={refreshComments}
                        />
                    ))}
                </div>
            )}
        </div>
    );
});

export default CommentsComponent;
