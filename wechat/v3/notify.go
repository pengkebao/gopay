package wechat

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pengkebao/gopay"
	"github.com/pengkebao/gopay/pkg/xlog"
)

type Resource struct {
	OriginalType   string `json:"original_type,omitempty"`
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

type V3DecryptResult struct {
	Appid           string             `json:"appid"`
	Mchid           string             `json:"mchid"`
	OutTradeNo      string             `json:"out_trade_no"`
	TransactionId   string             `json:"transaction_id"`
	TradeType       string             `json:"trade_type"`
	TradeState      string             `json:"trade_state"`
	TradeStateDesc  string             `json:"trade_state_desc"`
	BankType        string             `json:"bank_type"`
	Attach          string             `json:"attach"`
	SuccessTime     string             `json:"success_time"`
	Payer           *Payer             `json:"payer"`
	Amount          *Amount            `json:"amount"`
	SceneInfo       *SceneInfo         `json:"scene_info"`
	PromotionDetail []*PromotionDetail `json:"promotion_detail"`
}

type V3DecryptPartnerResult struct {
	SpAppid         string             `json:"sp_appid"`
	SpMchid         string             `json:"sp_mchid"`
	SubAppid        string             `json:"sub_appid"`
	SubMchid        string             `json:"sub_mchid"`
	OutTradeNo      string             `json:"out_trade_no"`
	TransactionId   string             `json:"transaction_id"`
	TradeType       string             `json:"trade_type"`
	TradeState      string             `json:"trade_state"`
	TradeStateDesc  string             `json:"trade_state_desc"`
	BankType        string             `json:"bank_type"`
	Attach          string             `json:"attach"`
	SuccessTime     string             `json:"success_time"`
	Payer           *PartnerPayer      `json:"payer"`
	Amount          *Amount            `json:"amount"`
	SceneInfo       *SceneInfo         `json:"scene_info"`
	PromotionDetail []*PromotionDetail `json:"promotion_detail"`
}

type V3DecryptRefundResult struct {
	Mchid               string        `json:"mchid"`
	OutTradeNo          string        `json:"out_trade_no"`
	TransactionId       string        `json:"transaction_id"`
	OutRefundNo         string        `json:"out_refund_no"`
	RefundId            string        `json:"refund_id"`
	RefundStatus        string        `json:"refund_status"`
	SuccessTime         string        `json:"success_time"`
	UserReceivedAccount string        `json:"user_received_account"`
	Amount              *RefundAmount `json:"amount"`
}

type V3DecryptPartnerRefundResult struct {
	SpMchid             string        `json:"sp_mchid"`
	SubMchid            string        `json:"sub_mchid"`
	OutTradeNo          string        `json:"out_trade_no"`
	TransactionId       string        `json:"transaction_id"`
	OutRefundNo         string        `json:"out_refund_no"`
	RefundId            string        `json:"refund_id"`
	RefundStatus        string        `json:"refund_status"`
	SuccessTime         string        `json:"success_time"`
	UserReceivedAccount string        `json:"user_received_account"`
	Amount              *RefundAmount `json:"amount"`
}

type V3DecryptCombineResult struct {
	CombineAppid      string       `json:"combine_appid"`
	CombineMchid      string       `json:"combine_mchid"`
	CombineOutTradeNo string       `json:"combine_out_trade_no"`
	SceneInfo         *SceneInfo   `json:"scene_info"`
	SubOrders         []*SubOrders `json:"sub_orders"`         // ???????????????????????????50
	CombinePayerInfo  *Payer       `json:"combine_payer_info"` // ???????????????
}

type V3DecryptScoreResult struct {
	Appid               string           `json:"appid"`
	Mchid               string           `json:"mchid"`
	OutOrderNo          string           `json:"out_order_no"`
	ServiceId           string           `json:"service_id"`
	Openid              string           `json:"openid"`
	State               string           `json:"state"`
	StateDescription    string           `json:"state_description"`
	TotalAmount         int              `json:"total_amount"`
	ServiceIntroduction string           `json:"service_introduction"`
	PostPayments        []*PostPayments  `json:"post_payments"`
	PostDiscounts       []*PostDiscounts `json:"post_discounts"`
	RiskFund            *RiskFund        `json:"risk_fund"`
	TimeRange           *TimeRange       `json:"time_range"`
	Location            *Location        `json:"location"`
	Attach              string           `json:"attach"`
	NotifyUrl           string           `json:"notify_url"`
	OrderId             string           `json:"order_id"`
	NeedCollection      bool             `json:"need_collection"`
	Collection          *Collection      `json:"collection"`
}

type V3DecryptProfitShareResult struct {
	SpMchid       string    `json:"sp_mchid"`       // ??????????????????
	SubMchid      string    `json:"sub_mchid"`      // ????????????
	TransactionId string    `json:"transaction_id"` // ???????????????
	OrderId       string    `json:"order_id"`       // ????????????/????????????
	OutOrderNo    string    `json:"out_order_no"`   // ????????????/????????????
	Receiver      *Receiver `json:"receiver"`
	SuccessTime   string    `json:"success_time"` // ????????????
}

type Receiver struct {
	Type        string `json:"type"`        // ?????????????????????
	Account     string `json:"account"`     // ?????????????????????
	Amount      int    `json:"amount"`      // ??????????????????
	Description string `json:"description"` // ??????/????????????
}

type V3DecryptBusifavorResult struct {
	EventType    string               `json:"event_type"`    // ????????????
	CouponCode   string               `json:"coupon_code"`   // ???code
	StockId      string               `json:"stock_id"`      // ?????????
	SendTime     string               `json:"send_time"`     // ????????????
	Openid       string               `json:"openid"`        // ????????????
	Unionid      string               `json:"unionid"`       // ??????????????????
	SendChannel  string               `json:"send_channel"`  // ????????????
	SendMerchant string               `json:"send_merchant"` // ???????????????
	AttachInfo   *BusifavorAttachInfo `json:"attach_info"`   // ??????????????????
}

type BusifavorAttachInfo struct {
	TransactionId   string `json:"transaction_id"`     // ??????????????????
	ActCode         string `json:"act_code"`           // ????????????????????????/???????????????ID
	HallCode        string `json:"hall_code"`          // ?????????ID
	HallBelongMchID int    `json:"hall_belong_mch_id"` // ????????????????????????
	CardID          string `json:"card_id"`            // ?????????ID
	Code            string `json:"code"`               // ?????????code
	ActivityID      string `json:"activity_id"`        // ????????????ID
}

type V3NotifyReq struct {
	Id           string    `json:"id"`
	CreateTime   string    `json:"create_time"`
	ResourceType string    `json:"resource_type"`
	EventType    string    `json:"event_type"`
	Summary      string    `json:"summary"`
	Resource     *Resource `json:"resource"`
	SignInfo     *SignInfo `json:"-"`
}

type V3NotifyRsp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// =====================================================================================================================

// ???????????????????????????????????? V3NotifyReq ?????????
func V3ParseNotify(req *http.Request) (notifyReq *V3NotifyReq, err error) {
	bs, err := ioutil.ReadAll(io.LimitReader(req.Body, int64(5<<20))) // default 5MB change the size you want;
	defer req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read request body error:%w", err)
	}
	si := &SignInfo{
		HeaderTimestamp: req.Header.Get(HeaderTimestamp),
		HeaderNonce:     req.Header.Get(HeaderNonce),
		HeaderSignature: req.Header.Get(HeaderSignature),
		HeaderSerial:    req.Header.Get(HeaderSerial),
		SignBody:        string(bs),
	}
	notifyReq = &V3NotifyReq{SignInfo: si}
	if err = json.Unmarshal(bs, notifyReq); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(%s, %+v)???%w", string(bs), notifyReq, err)
	}
	return notifyReq, nil
}

// Deprecated
// ???????????? VerifySignByPK()
func (v *V3NotifyReq) VerifySign(wxPkContent string) (err error) {
	if v.SignInfo != nil {
		return V3VerifySign(v.SignInfo.HeaderTimestamp, v.SignInfo.HeaderNonce, v.SignInfo.SignBody, v.SignInfo.HeaderSignature, wxPkContent)
	}
	return errors.New("verify notify sign, bug SignInfo is nil")
}

// ??????????????????
// wxPublicKey?????????????????????????????????????????? client.WxPublicKeyMap() ????????????????????? signInfo.HeaderSerial ?????????????????????
// ???????????? VerifySignByPKMap()
func (v *V3NotifyReq) VerifySignByPK(wxPublicKey *rsa.PublicKey) (err error) {
	if v.SignInfo != nil {
		return V3VerifySignByPK(v.SignInfo.HeaderTimestamp, v.SignInfo.HeaderNonce, v.SignInfo.SignBody, v.SignInfo.HeaderSignature, wxPublicKey)
	}
	return errors.New("verify notify sign, bug SignInfo is nil")
}

// ??????????????????
// wxPublicKey?????????????????????????????????????????? client.WxPublicKeyMap() ??????
func (v *V3NotifyReq) VerifySignByPKMap(wxPublicKeyMap map[string]*rsa.PublicKey) (err error) {
	if v.SignInfo != nil && wxPublicKeyMap != nil {
		return V3VerifySignByPK(v.SignInfo.HeaderTimestamp, v.SignInfo.HeaderNonce, v.SignInfo.SignBody, v.SignInfo.HeaderSignature, wxPublicKeyMap[v.SignInfo.HeaderSerial])
	}
	return errors.New("verify notify sign, bug SignInfo or wxPublicKeyMap is nil")
}

// ?????? ???????????? ????????????????????????
func (v *V3NotifyReq) DecryptCipherText(apiV3Key string) (result *V3DecryptResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ?????? ??????????????? ????????????????????????
func (v *V3NotifyReq) DecryptPartnerCipherText(apiV3Key string) (result *V3DecryptPartnerResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptPartnerNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ?????? ???????????? ????????????????????????
func (v *V3NotifyReq) DecryptRefundCipherText(apiV3Key string) (result *V3DecryptRefundResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptRefundNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ?????? ??????????????? ????????????????????????
func (v *V3NotifyReq) DecryptPartnerRefundCipherText(apiV3Key string) (result *V3DecryptPartnerRefundResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptPartnerRefundNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ?????? ???????????? ????????????????????????
func (v *V3NotifyReq) DecryptCombineCipherText(apiV3Key string) (result *V3DecryptCombineResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptCombineNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ?????? ????????? ????????????????????????
func (v *V3NotifyReq) DecryptScoreCipherText(apiV3Key string) (result *V3DecryptScoreResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptScoreNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ??????????????????????????????????????????
func (v *V3NotifyReq) DecryptProfitShareCipherText(apiV3Key string) (result *V3DecryptProfitShareResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptProfitShareNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// ???????????????????????????????????????
func (v *V3NotifyReq) DecryptBusifavorCipherText(apiV3Key string) (result *V3DecryptBusifavorResult, err error) {
	if v.Resource != nil {
		result, err = V3DecryptBusifavorNotifyCipherText(v.Resource.Ciphertext, v.Resource.Nonce, v.Resource.AssociatedData, apiV3Key)
		if err != nil {
			bytes, _ := json.Marshal(v)
			return nil, fmt.Errorf("V3NotifyReq(%s) decrypt cipher text error(%w)", string(bytes), err)
		}
		return result, nil
	}
	return nil, errors.New("notify data Resource is nil")
}

// Deprecated
// ???????????????????????????????????? wechat.V3ParseNotify()
// ???????????????????????????????????? gopay.BodyMap
func V3ParseNotifyToBodyMap(req *http.Request) (bm gopay.BodyMap, err error) {
	bs, err := ioutil.ReadAll(io.LimitReader(req.Body, int64(3<<20))) // default 3MB change the size you want;
	defer req.Body.Close()
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	bm = make(gopay.BodyMap)
	if err = json.Unmarshal(bs, &bm); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(%s)???%w", string(bs), err)
	}
	return bm, nil
}
