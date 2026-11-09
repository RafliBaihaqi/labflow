package repo

import (
	"labflow/models"
	"labflow/sqlc"
)

// The services work with the domain models in package models. These functions
// convert the rows sqlc generates into those models.

func toUser(u sqlc.User) models.User {
	return models.User{
		ID:                 u.ID,
		XID:                u.Xid,
		Username:           u.Username,
		Password:           u.Password,
		Email:              u.Email,
		PhoneNumber:        deref(u.PhoneNumber),
		Name:               u.Name,
		Company:            deref(u.Company),
		Address:            deref(u.Address),
		VerificationToken:  u.VerificationToken,
		VerificationExp:    u.VerificationExp,
		IsVerified:         u.IsVerified,
		ResetPasswordToken: u.ResetPasswordToken,
		ResetPasswordExp:   u.ResetPasswordExp,
		SessionID:          u.SessionID,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		UpdatedBy:          u.UpdatedBy,
		Avatar:             u.Avatar,
	}
}

func toReservation(r sqlc.Reservation) models.Reservation {
	return models.Reservation{
		ID:                  r.ID,
		XID:                 r.Xid,
		UserID:              r.UserID,
		ReservationDate:     r.ReservationDate,
		Status:              r.Status,
		SampleName:          r.SampleName,
		SampleAmount:        r.SampleAmount,
		ReceivedNote:        r.ReceivedNote,
		TestRequestNumber:   r.TestRequestNumber,
		SampleResidue:       r.SampleResidue,
		SpecialInstructions: r.SpecialInstructions,
		POURL:               r.PoUrl,
		PONumber:            r.PoNumber,
		SampleType:          r.SampleType,
		FeedbackRequired:    r.FeedbackRequired,
		MRApproved:          r.MrApproved,
		HeadApproved:        r.HeadApproved,
		MRApprovalDate:      r.MrApprovalDate,
		HeadApprovalDate:    r.HeadApprovalDate,
		SamplePhoto:         r.SamplePhoto,
		Note:                r.Note,
		StartTestDate:       r.StartTestDate,
		EndTestDate:         r.EndTestDate,
		SAPNotes:            r.SapNotes,
		BatchNumber:         r.BatchNumber,
		LeadTime:            r.LeadTime,
		SampleTakenDate:     r.SampleTakenDate,
		SampleTakenBy:       r.SampleTakenBy,
		IsKan:               r.IsKan,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
		UpdateBy:            r.UpdatedBy,
	}
}

func toTest(t sqlc.Test) models.Test {
	return models.Test{
		ID:            t.ID,
		XID:           t.Xid,
		ReservationID: t.ReservationID,
		TestListID:    t.TestListID,
		Method:        t.Method,
		Criteria:      t.Criteria,
		Status:        t.Status,
		TestUnit:      t.TestUnit,
		TestResult:    t.TestResult,
		TestSummary:   t.TestSummary,
		ApprovedDate:  t.ApprovedDate,
		Notes:         t.Notes,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		UpdateBy:      t.UpdatedBy,
		IsApproved:    t.IsApproved,
	}
}

func toTestList(tl sqlc.TestList) models.TestList {
	return models.TestList{
		ID:                 tl.ID,
		XID:                tl.Xid,
		CategoryID:         tl.CategoryID,
		Value:              tl.Value,
		Label:              tl.Label,
		SampleAmountLiquid: tl.SampleAmountLiquid,
		SampleAmountSolid:  tl.SampleAmountSolid,
		SampleAmountPowder: tl.SampleAmountPowder,
		Leadtime:           tl.Leadtime,
		Version:            tl.Version,
		IsCurrent:          tl.IsCurrent,
		CreatedAt:          tl.CreatedAt,
		UpdatedAt:          tl.UpdatedAt,
	}
}

func toCategory(c sqlc.Category) models.Category {
	return models.Category{
		ID:        c.ID,
		XID:       c.Xid,
		Name:      c.Name,
		Label:     c.Label,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toTestListWithCategory(tl sqlc.TestList, c sqlc.Category) models.TestList {
	out := toTestList(tl)
	out.Category = toCategory(c)
	return out
}

func toQuotation(q sqlc.Quotation) models.Quotation {
	return models.Quotation{
		ID:          q.ID,
		XID:         q.Xid,
		Email:       q.Email,
		CompanyName: q.CompanyName,
		PhoneNumber: q.PhoneNumber,
		Message:     q.Message,
		Status:      q.Status,
		CreatedAt:   q.CreatedAt,
		UpdatedAt:   q.UpdatedAt,
	}
}
