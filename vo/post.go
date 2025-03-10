package vo

type CreatePostRequest struct {
	CategoryId uint   `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,max=50"`
	HeadImgUrl string `json:"head_img_url"`
	Content    string `json:"content" binding:"required"`
}
