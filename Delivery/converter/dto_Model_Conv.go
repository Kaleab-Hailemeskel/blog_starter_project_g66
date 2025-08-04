package conv

import (
	"blog_starter_project_g66/Domain"
)

func ChangeToDTOUserfunc(domainUser *domain.User) *domain.UserDTO {
	return &domain.UserDTO{
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
func ChangeToDomainUser(udto *domain.UserDTO) *domain.User {
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

func ChangeToDomainBlog(bdto *domain.BlogDTO) *domain.Blog {
	return &domain.Blog{
		Title:       bdto.Title,
		Tags:        bdto.Tags,
		Author:      bdto.Author,
		Description: bdto.Description,
		LastUpdate:  bdto.LastUpdate,
	}
}
func ChangeToDTOBlog(bdto *domain.Blog) *domain.BlogDTO {
	return &domain.BlogDTO{
		Title:       bdto.Title,
		Tags:        bdto.Tags,
		Author:      bdto.Author,
		Description: bdto.Description,
		LastUpdate:  bdto.LastUpdate,
	}
}
func ChangeToDomainPopularity(pdto *domain.PopularityDTO) *domain.Popularity {
	return &domain.Popularity{
		ViewCount: pdto.ViewCount,
		Likes:     pdto.Likes,
		Dislikes:  pdto.Dislikes,
		Comments:  changeToListDomainComment(pdto.Comments),
	}
}

func ChangeToDTOPopularity(pdto *domain.Popularity) *domain.PopularityDTO {
	return &domain.PopularityDTO{
		ViewCount: pdto.ViewCount,
		Likes:     pdto.Likes,
		Dislikes:  pdto.Dislikes,
		Comments:  changeToListDTOComment(pdto.Comments),
	}
}
func ChangeToDomainComment(cdto *domain.CommentDTO) *domain.Comment {
	return &domain.Comment{
		UserName: cdto.UserName,
		Comment:  cdto.Comment,
	}
}
func ChangeToDTOComment(cdto *domain.Comment) *domain.CommentDTO {
	return &domain.CommentDTO{
		UserName: cdto.UserName,
		Comment:  cdto.Comment,
	}
}

func changeToListDomainComment(lctdo []*domain.CommentDTO) []*domain.Comment {
	var listDomainComment []*domain.Comment
	for _, val := range lctdo {
		listDomainComment = append(listDomainComment, ChangeToDomainComment(val))
	}
	return listDomainComment
}
func changeToListDTOComment(lctdo []*domain.Comment) []*domain.CommentDTO {
	var listDTOComment []*domain.CommentDTO
	for _, val := range lctdo {
		listDTOComment = append(listDTOComment, ChangeToDTOComment(val))
	}
	return listDTOComment
}
func ChangeToDomainVerification(udto *domain.UserUnverifiedDTO) *domain.UserUnverified{
    return &domain.UserUnverified{
		UserName: udto.UserName,
        Email:     udto.Email,
        OTP:       udto.OTP,
		Password: udto.Password,
		Role: udto.Role,
        ExpiresAt: udto.ExpiresAt,
    }
}

func ChangeUnverfiedToVerified(u *domain.UserUnverifiedDTO) *domain.User{
	return&domain.User{
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		Role: u.Role,
	}
}