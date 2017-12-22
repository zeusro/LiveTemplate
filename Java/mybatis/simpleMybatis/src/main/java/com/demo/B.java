package com.demo;

public class B {

    private int id;
    private String a;
    private long b;

    public B() {
        super();
    }


    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }


    public String getA() {
        return a;
    }

    public void setA(String a) {
        this.a = a;
    }

    public long getB() {
        return b;
    }

    public void setB(long b) {
        this.b = b;
    }

    @Override
    public String toString() {
        //http://blog.csdn.net/lonely_fireworks/article/details/7962171
        return String.format("id:%s;a:%s;b:%s;", id, a, b);
    }
}

