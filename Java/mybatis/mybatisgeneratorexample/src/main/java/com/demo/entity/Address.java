package com.demo.entity;

import java.util.Date;

public class Address {
    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.address_id
     *
     * @mbg.generated
     */
    private Short addressId;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.address
     *
     * @mbg.generated
     */
    private String address;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.address2
     *
     * @mbg.generated
     */
    private String address2;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.district
     *
     * @mbg.generated
     */
    private String district;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.city_id
     *
     * @mbg.generated
     */
    private Short cityId;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.postal_code
     *
     * @mbg.generated
     */
    private String postalCode;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.phone
     *
     * @mbg.generated
     */
    private String phone;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.last_update
     *
     * @mbg.generated
     */
    private Date lastUpdate;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column address.location
     *
     * @mbg.generated
     */
    private byte[] location;

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table address
     *
     * @mbg.generated
     */
    public Address(Short addressId, String address, String address2, String district, Short cityId, String postalCode, String phone, Date lastUpdate, byte[] location) {
        this.addressId = addressId;
        this.address = address;
        this.address2 = address2;
        this.district = district;
        this.cityId = cityId;
        this.postalCode = postalCode;
        this.phone = phone;
        this.lastUpdate = lastUpdate;
        this.location = location;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.address_id
     *
     * @return the value of address.address_id
     *
     * @mbg.generated
     */
    public Short getAddressId() {
        return addressId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.address
     *
     * @return the value of address.address
     *
     * @mbg.generated
     */
    public String getAddress() {
        return address;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.address2
     *
     * @return the value of address.address2
     *
     * @mbg.generated
     */
    public String getAddress2() {
        return address2;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.district
     *
     * @return the value of address.district
     *
     * @mbg.generated
     */
    public String getDistrict() {
        return district;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.city_id
     *
     * @return the value of address.city_id
     *
     * @mbg.generated
     */
    public Short getCityId() {
        return cityId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.postal_code
     *
     * @return the value of address.postal_code
     *
     * @mbg.generated
     */
    public String getPostalCode() {
        return postalCode;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.phone
     *
     * @return the value of address.phone
     *
     * @mbg.generated
     */
    public String getPhone() {
        return phone;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.last_update
     *
     * @return the value of address.last_update
     *
     * @mbg.generated
     */
    public Date getLastUpdate() {
        return lastUpdate;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column address.location
     *
     * @return the value of address.location
     *
     * @mbg.generated
     */
    public byte[] getLocation() {
        return location;
    }
}