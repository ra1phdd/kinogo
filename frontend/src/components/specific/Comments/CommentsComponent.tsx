import React, { useState, useCallback } from 'react';
import useComments from '@/hooks/useComments.ts';
import Editor from '@components/specific/Editor.tsx';
import Comment from './Comment';
import { Comments } from '@components/gRPC.tsx';
import {useAuth} from "@/contexts/Auth.tsx";

const CommentsComponent: React.FC<{ id: number }> = ({ id }) => {
    const { comments, loading } = useComments(id);
    const { isAuthenticated  } = useAuth();
    const [replyToParentId, setReplyToParentId] = useState<number | null>(null);
    const [replyToParentComment, setReplyToParentComment] = useState<Comments | null>(null);
    const [editComment, setEditComment] = useState<Comments | null>(null);

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

    return (
        <>
            {isAuthenticated && (
                <Editor
                    parentId={replyToParentId}
                    parentComment={replyToParentComment}
                    movieId={id}
                    onSubmit={handleSubmit}
                    onCancelReply={handleCancelReply}
                    editComment={editComment}
                />
            )}
            {comments.map((comment) => (
                <Comment key={comment.id} comment={comment} onReply={handleReply} onEdit={handleEdit} />
            ))}
            {loading && <p>Загрузка...</p>}
            {!loading && comments.length === 0 && <p>No comments found.</p>}
        </>
    );
};

export default CommentsComponent;
