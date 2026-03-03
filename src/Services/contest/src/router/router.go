package router

import (
	"contest/src/components"
	contestcontroller "contest/src/controller/contest"
	contestsubmissioncontroller "contest/src/controller/contest-submission"
	contestantcontroller "contest/src/controller/contestant"
	ratingcontroller "contest/src/controller/rating"
	scoreboardcontroller "contest/src/controller/scoreboard"
	"contest/src/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, appContext components.AppContext) {
	contestInteractor := appContext.GetContestInteractor()
	contestantInteractor := appContext.GetContestantInteractor()
	contestsubmissionInteractor := appContext.GetContestSubmissionInteractor()
	scoreboardInteractor := appContext.GetScoreboardInteractor()
	ratingInteractor := appContext.GetRatingInteractor()
	contest := group.Group("/contest")
	{
		auth := contest.Group("")
		auth.Use(middleware.OptionalAuth())
		auth.GET("", contestcontroller.GetAllContest(contestInteractor))
		auth.GET("/:contest_id", contestcontroller.GetContest(contestInteractor))

		auth.POST("/scoreboard", scoreboardcontroller.BuildScoreboardSnapshot(scoreboardInteractor))
		auth.POST("/scoreboard/at", scoreboardcontroller.GetScoreboardAt(scoreboardInteractor))

		auth.GET("/:contest_id/rating", ratingcontroller.GetContestRatingResults(ratingInteractor))
	}
	{
		auth := contest.Group("")
		auth.Use(middleware.RequireAuth())
		auth.POST("/create", contestcontroller.Create(contestInteractor))
		auth.POST("/edit", contestcontroller.Edit(contestInteractor))
		auth.PATCH("/:contest_id/patch", contestcontroller.Patch(contestInteractor))
		auth.PUT("/:contest_id/manage/problems", contestcontroller.ManageProblems(contestInteractor))

		auth.POST("/submit", contestantcontroller.Submit(contestantInteractor))
		auth.POST("/register", contestantcontroller.Register(contestantInteractor))
		auth.POST("/unregister", contestantcontroller.Unregister(contestantInteractor))

		auth.POST("/:contest_id/compute_rating", ratingcontroller.ComputeContestRatings(ratingInteractor))
	}
	{
		internal := contest.Group("/internal")
		internal.Use(middleware.Internal())
		internal.POST("/ingest-submission", contestsubmissioncontroller.InternalIngestContestSubmission(contestsubmissionInteractor))
		internal.GET("/problem-lock/:problem_id", contestcontroller.InternalProblemLock(contestInteractor))
	}

	rating := group.Group("/rating")
	{
		rating.GET("/user/:username", ratingcontroller.GetUserRating(ratingInteractor))
		rating.GET("/user/:username/history", ratingcontroller.GetUserRatingHistory(ratingInteractor))
	}
}
