package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r *CommentRepository) Create(ctx context.Context, c *domain.Comment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO comments (id, entity_type, entity_id, user_id, parent_id, content, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.EntityType, c.EntityID, c.UserID, c.ParentID, c.Content, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	c := &domain.Comment{}
	err := r.pool.QueryRow(ctx,
		`SELECT c.id, c.entity_type, c.entity_id, c.user_id, c.parent_id, c.content,
		        c.created_at, c.updated_at, u.name, u.email
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.id = $1`,
		id,
	).Scan(&c.ID, &c.EntityType, &c.EntityID, &c.UserID, &c.ParentID, &c.Content,
		&c.CreatedAt, &c.UpdatedAt, &c.UserName, &c.UserEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return c, err
}

func (r *CommentRepository) ListByEntity(ctx context.Context, entityType domain.CommentEntityType, entityID uuid.UUID) ([]domain.Comment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT c.id, c.entity_type, c.entity_id, c.user_id, c.parent_id, c.content,
		        c.created_at, c.updated_at, u.name, u.email
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.entity_type = $1 AND c.entity_id = $2 AND c.parent_id IS NULL
		 ORDER BY c.created_at ASC`,
		entityType, entityID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topLevel []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.EntityType, &c.EntityID, &c.UserID, &c.ParentID, &c.Content,
			&c.CreatedAt, &c.UpdatedAt, &c.UserName, &c.UserEmail); err != nil {
			return nil, err
		}
		topLevel = append(topLevel, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range topLevel {
		replies, err := r.loadReplies(ctx, topLevel[i].ID)
		if err != nil {
			return nil, err
		}
		topLevel[i].Replies = replies
	}
	return topLevel, nil
}

func (r *CommentRepository) loadReplies(ctx context.Context, parentID uuid.UUID) ([]domain.Comment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT c.id, c.entity_type, c.entity_id, c.user_id, c.parent_id, c.content,
		        c.created_at, c.updated_at, u.name, u.email
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.parent_id = $1
		 ORDER BY c.created_at ASC`,
		parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.EntityType, &c.EntityID, &c.UserID, &c.ParentID, &c.Content,
			&c.CreatedAt, &c.UpdatedAt, &c.UserName, &c.UserEmail); err != nil {
			return nil, err
		}
		replies = append(replies, c)
	}
	return replies, rows.Err()
}

func (r *CommentRepository) Update(ctx context.Context, c *domain.Comment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE comments SET content = $2, updated_at = $3 WHERE id = $1`,
		c.ID, c.Content, c.UpdatedAt,
	)
	return err
}

func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	return err
}
