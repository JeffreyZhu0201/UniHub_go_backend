package DTO

import "time"

type CreateDingRequest struct {
	Title     string    `json:"title" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	//Type       string    `json:"type" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
	Latitude   float64   `json:"latitude" binding:"required"`
	Longitude  float64   `json:"longitude" binding:"required"`
	Radius     uint      `json:"radius" binding:"required"`
	StudentId  uint      `json:"student_id"`
	DeptId     uint      `json:"dept_id"`
	ClassId    uint      `json:"class_id"`
	LauncherId uint      `json:"launcher_id"`
}
