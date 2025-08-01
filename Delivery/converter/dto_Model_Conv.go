package conv

import (
	"blog_starter_project_g66/Delivery/controllers"
	domain "blog_starter_project_g66/Domain"
)

func ChangeToDTOUserfunc(domainUser *domain.User) *controllers.UserDTO {
	return &controllers.UserDTO{
		UserName:       domainUser.UserName,
		PersonalBio:    domainUser.PersonalBio,
		ProfilePic:     domainUser.ProfilePic,
		Email:          domainUser.Email,
		PhoneNum:       domainUser.PhoneNum,
		TelegramHandle: domainUser.TelegramHandle,
		Password:       domainUser.Password,
		Role:           domainUser.Role,
	}
}
func ChangeToDomainUser(udto *controllers.UserDTO) *domain.User {
	return &domain.User{
		UserName:       udto.UserName,
		PersonalBio:    udto.PersonalBio,
		ProfilePic:     udto.ProfilePic,
		Email:          udto.Email,
		PhoneNum:       udto.PhoneNum,
		TelegramHandle: udto.TelegramHandle,
		Password:       udto.Password,
		Role:           udto.Role,
	}
}

func ChangeToDomainBlog(bdto *controllers.BlogDTO) *domain.Blog {
	return &domain.Blog{
		Title:       bdto.Title,
		Tags:        bdto.Tags,
		Author:      bdto.Author,
		Description: bdto.Description,
		LastUpdate:  bdto.LastUpdate,
	}
}
func ChangeToDTOBlog(bdto *domain.Blog) *controllers.BlogDTO {
	return &controllers.BlogDTO{
		Title:       bdto.Title,
		Tags:        bdto.Tags,
		Author:      bdto.Author,
		Description: bdto.Description,
		LastUpdate:  bdto.LastUpdate,
	}
}
func ChangeToDomainPopularity(pdto *controllers.PopularityDTO) *domain.Popularity {
	return &domain.Popularity{
		ViewCount: pdto.ViewCount,
		Likes:     pdto.Likes,
		Dislikes:  pdto.Dislikes,
		Comments:  changeToListDomainComment(pdto.Comments),
	}
}

func ChangeToDTOPopularity(pdto *domain.Popularity) *controllers.PopularityDTO {
	return &controllers.PopularityDTO{
		ViewCount: pdto.ViewCount,
		Likes:     pdto.Likes,
		Dislikes:  pdto.Dislikes,
		Comments:  changeToListDTOComment(pdto.Comments),
	}
}
func ChangeToDomainComment(cdto *controllers.CommentDTO) *domain.Comment {
	return &domain.Comment{
		UserName: cdto.UserName,
		Comment:  cdto.Comment,
	}
}
func ChangeToDTOComment(cdto *domain.Comment) *controllers.CommentDTO {
	return &controllers.CommentDTO{
		UserName: cdto.UserName,
		Comment:  cdto.Comment,
	}
}

func changeToListDomainComment(lctdo []*controllers.CommentDTO) []*domain.Comment {
	var listDomainComment []*domain.Comment
	for _, val := range lctdo {
		listDomainComment = append(listDomainComment, ChangeToDomainComment(val))
	}
	return listDomainComment
}
func changeToListDTOComment(lctdo []*domain.Comment) []*controllers.CommentDTO {
	var listDTOComment []*controllers.CommentDTO
	for _, val := range lctdo {
		listDTOComment = append(listDTOComment, ChangeToDTOComment(val))
	}
	return listDTOComment
}
