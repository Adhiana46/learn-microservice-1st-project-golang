import React from "react";

const CommentList = ({ comments }) => {
  const renderedComments = comments.map((comment) => {
    let content = comment.content;
    switch (comment.status) {
      case "rejected":
        content = "rejected";
        break;
      case "approved":
        content = comment.content;
        break;
      default:
        content = "** Awaiting moderation **";
    }

    return <li key={comment.id}>{content}</li>;
  });

  return <ul>{renderedComments}</ul>;
};

export default CommentList;
