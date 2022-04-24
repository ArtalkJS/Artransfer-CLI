package waline

import "github.com/ArtalkJS/Artransfer-CLI/lib"

type WalineCmd struct {
}

func (cmd *WalineCmd) Run(ctx *lib.Context) error {
	ctx.SrcType = "waline"

	var wComments []Comment
	ctx.DB.Find(&wComments)

	artrans := WalineToArtrans(wComments)
	json := ctx.ArtransToJson(artrans)

	ctx.Export(json)

	return nil
}

func WalineToArtrans(wComments []Comment) []lib.Artran {
	artrans := []lib.Artran{}

	for _, wC := range wComments {
		artran := lib.Artran{
			ID:        lib.ToString(wC.ID),
			Rid:       lib.ToString(wC.Rid),
			Content:   wC.Comment,
			UA:        wC.Ua,
			IP:        wC.Ip,
			IsPending: lib.ToString(wC.Status == "waiting"),
			CreatedAt: wC.InsertedAt.String(),
			UpdatedAt: wC.UpdatedAt.String(),
			Nick:      wC.Nick,
			Email:     wC.Mail,
			Link:      wC.Link,
			PageKey:   wC.Url,
		}

		artrans = append(artrans, artran)
	}

	return artrans
}
